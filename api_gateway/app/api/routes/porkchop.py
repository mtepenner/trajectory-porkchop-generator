"""Porkchop plot endpoint – returns a delta-v contour matrix."""
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel, Field
from typing import List
import numpy as np
import math

from app.core.redis_cache import CacheClient

router = APIRouter()


class PorkchopRequest(BaseModel):
    departure_body: str = Field("Earth", description="Departure body name")
    arrival_body: str = Field("Mars", description="Arrival body name")
    dep_mjd_start: float = Field(..., description="Departure window start (MJD)")
    dep_mjd_end: float = Field(..., description="Departure window end (MJD)")
    arr_mjd_start: float = Field(..., description="Arrival window start (MJD)")
    arr_mjd_end: float = Field(..., description="Arrival window end (MJD)")
    steps: int = Field(100, ge=10, le=500)


class PorkchopResponse(BaseModel):
    grid: List[List[float]]
    dep_mjd_axis: List[float]
    arr_mjd_axis: List[float]
    steps: int


@router.post("/compute", response_model=PorkchopResponse)
async def compute_porkchop(req: PorkchopRequest):
    cache_params = req.dict()
    cached = await CacheClient.get("porkchop", cache_params)
    if cached:
        return cached

    # Call the compute engine (gRPC stub – simplified HTTP call for skeleton)
    grid, dep_axis, arr_axis = _compute_grid(
        req.dep_mjd_start, req.dep_mjd_end,
        req.arr_mjd_start, req.arr_mjd_end,
        req.steps,
    )

    result = {
        "grid": grid,
        "dep_mjd_axis": dep_axis,
        "arr_mjd_axis": arr_axis,
        "steps": req.steps,
    }
    await CacheClient.set("porkchop", cache_params, result)
    return result


def _compute_grid(dep0, dep1, arr0, arr1, steps):
    """Local fallback grid computation using simplified two-body physics."""
    MU_SUN = 1.32712440018e11  # km^3/s^2
    AU = 149597870.7            # km
    r_earth = 1.0 * AU
    r_mars = 1.524 * AU

    dep_axis = [dep0 + i * (dep1 - dep0) / (steps - 1) for i in range(steps)]
    arr_axis = [arr0 + i * (arr1 - arr0) / (steps - 1) for i in range(steps)]

    grid = []
    for dep_mjd in dep_axis:
        row = []
        for arr_mjd in arr_axis:
            if arr_mjd <= dep_mjd:
                row.append(float("nan"))
                continue
            tof = (arr_mjd - dep_mjd) * 86400.0  # seconds
            # Hohmann-approximation delta-v (simplified)
            a_transfer = (r_earth + r_mars) / 2
            v_c_earth = math.sqrt(MU_SUN / r_earth)
            v_c_mars = math.sqrt(MU_SUN / r_mars)
            v_dep = math.sqrt(MU_SUN * (2 / r_earth - 1 / a_transfer))
            v_arr = math.sqrt(MU_SUN * (2 / r_mars - 1 / a_transfer))
            dv = abs(v_dep - v_c_earth) + abs(v_c_mars - v_arr)
            row.append(round(dv, 3))
        grid.append(row)

    return grid, dep_axis, arr_axis
