"""Trajectory 3-D flight-path endpoint."""
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel, Field
from typing import List
import math

from app.core.redis_cache import CacheClient

router = APIRouter()


class TrajectoryRequest(BaseModel):
    departure_body: str = Field("Earth")
    arrival_body: str = Field("Mars")
    dep_mjd: float = Field(..., description="Departure date (MJD)")
    arr_mjd: float = Field(..., description="Arrival date (MJD)")
    num_points: int = Field(200, ge=10, le=1000)


class TrajectoryPoint(BaseModel):
    t: float    # time from departure (s)
    x: float    # km
    y: float    # km
    z: float    # km


class TrajectoryResponse(BaseModel):
    path: List[TrajectoryPoint]
    dep_mjd: float
    arr_mjd: float
    tof_days: float
    delta_v_total: float


@router.post("/compute", response_model=TrajectoryResponse)
async def compute_trajectory(req: TrajectoryRequest):
    cache_params = req.dict()
    cached = await CacheClient.get("trajectory", cache_params)
    if cached:
        return cached

    if req.arr_mjd <= req.dep_mjd:
        raise HTTPException(status_code=400, detail="Arrival must be after departure")

    path = _propagate(req.dep_mjd, req.arr_mjd, req.num_points)

    result = {
        "path": path,
        "dep_mjd": req.dep_mjd,
        "arr_mjd": req.arr_mjd,
        "tof_days": req.arr_mjd - req.dep_mjd,
        "delta_v_total": 5.6,  # approximate placeholder km/s
    }
    await CacheClient.set("trajectory", cache_params, result)
    return result


def _propagate(dep_mjd, arr_mjd, n):
    """Generate a simplified heliocentric elliptical arc between Earth and Mars."""
    AU = 149597870.7
    r_earth = 1.0 * AU
    r_mars = 1.524 * AU
    a = (r_earth + r_mars) / 2
    tof = (arr_mjd - dep_mjd) * 86400.0

    points = []
    for i in range(n):
        frac = i / (n - 1)
        theta = math.pi * frac  # half-ellipse
        r = a * (1 - 0.2 * math.cos(theta))
        points.append({
            "t": frac * tof,
            "x": r * math.cos(theta),
            "y": r * math.sin(theta),
            "z": 0.0,
        })
    return points
