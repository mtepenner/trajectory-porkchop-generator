from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel
from typing import Optional
import asyncio

from app.api.routes import porkchop, trajectory
from app.core.redis_cache import CacheClient

app = FastAPI(
    title="Trajectory Porkchop Generator API",
    description="Interplanetary trajectory computation and caching service",
    version="1.0.0",
)

app.include_router(porkchop.router, prefix="/api/v1/porkchop", tags=["porkchop"])
app.include_router(trajectory.router, prefix="/api/v1/trajectory", tags=["trajectory"])


@app.on_event("startup")
async def startup_event():
    await CacheClient.initialize()


@app.on_event("shutdown")
async def shutdown_event():
    await CacheClient.close()


@app.get("/health")
async def health():
    return {"status": "ok"}
