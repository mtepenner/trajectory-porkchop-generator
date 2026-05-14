"""Redis cache client for trajectory matrix results."""
import json
import hashlib
from typing import Optional, Any

import redis.asyncio as redis

_client: Optional[redis.Redis] = None
CACHE_TTL = 3600  # 1 hour


class CacheClient:
    @classmethod
    async def initialize(cls, url: str = "redis://redis:6379"):
        global _client
        _client = redis.from_url(url, encoding="utf-8", decode_responses=True)

    @classmethod
    async def close(cls):
        global _client
        if _client:
            await _client.close()

    @classmethod
    def _make_key(cls, namespace: str, params: dict) -> str:
        raw = json.dumps(params, sort_keys=True)
        digest = hashlib.sha256(raw.encode()).hexdigest()[:16]
        return f"{namespace}:{digest}"

    @classmethod
    async def get(cls, namespace: str, params: dict) -> Optional[Any]:
        if not _client:
            return None
        key = cls._make_key(namespace, params)
        data = await _client.get(key)
        if data:
            return json.loads(data)
        return None

    @classmethod
    async def set(cls, namespace: str, params: dict, value: Any):
        if not _client:
            return
        key = cls._make_key(namespace, params)
        await _client.setex(key, CACHE_TTL, json.dumps(value))
