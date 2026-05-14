import { useState, useCallback } from 'react';
import axios from 'axios';

const API_BASE = process.env.REACT_APP_API_URL ?? 'http://localhost:8000';

export interface GridData {
  grid: number[][];
  dep_mjd_axis: number[];
  arr_mjd_axis: number[];
  steps: number;
}

export interface TrajectoryPoint {
  t: number;
  x: number;
  y: number;
  z: number;
}

export interface Metrics {
  delta_v_total: number;
  tof_days: number;
  dep_mjd: number;
  arr_mjd: number;
}

export function useTrajectoryData() {
  const [gridData, setGridData] = useState<GridData | null>(null);
  const [trajectoryPath, setTrajectoryPath] = useState<TrajectoryPoint[] | null>(null);
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchGrid = useCallback(
    async (depStart: number, depEnd: number, arrStart: number, arrEnd: number, steps: number) => {
      setLoading(true);
      setError(null);
      try {
        const res = await axios.post(`${API_BASE}/api/v1/porkchop/compute`, {
          departure_body: 'Earth',
          arrival_body: 'Mars',
          dep_mjd_start: depStart,
          dep_mjd_end: depEnd,
          arr_mjd_start: arrStart,
          arr_mjd_end: arrEnd,
          steps,
        });
        setGridData(res.data);
      } catch (e: any) {
        setError(e.message ?? 'Failed to fetch grid');
      } finally {
        setLoading(false);
      }
    },
    []
  );

  const fetchPath = useCallback(async (dep_mjd: number, arr_mjd: number) => {
    setLoading(true);
    setError(null);
    try {
      const res = await axios.post(`${API_BASE}/api/v1/trajectory/compute`, {
        departure_body: 'Earth',
        arrival_body: 'Mars',
        dep_mjd,
        arr_mjd,
        num_points: 200,
      });
      setTrajectoryPath(res.data.path);
      setMetrics({
        delta_v_total: res.data.delta_v_total,
        tof_days: res.data.tof_days,
        dep_mjd: res.data.dep_mjd,
        arr_mjd: res.data.arr_mjd,
      });
    } catch (e: any) {
      setError(e.message ?? 'Failed to fetch trajectory');
    } finally {
      setLoading(false);
    }
  }, []);

  return { gridData, trajectoryPath, metrics, loading, error, fetchGrid, fetchPath };
}
