import React from 'react';
import { Metrics } from '../hooks/useTrajectoryData';

interface Props { metrics: Metrics | null; }

const DeltaVMetrics: React.FC<Props> = ({ metrics }) => {
  if (!metrics) {
    return (
      <div style={{ background: '#111', borderRadius: 8, padding: '1rem', minWidth: 220 }}>
        <h3 style={{ color: '#7af', margin: 0 }}>ΔV Metrics</h3>
        <p style={{ color: '#888' }}>Click a cell to compute a trajectory.</p>
      </div>
    );
  }

  const AU = 149597870.7;

  return (
    <div style={{ background: '#111', borderRadius: 8, padding: '1rem', minWidth: 220 }}>
      <h3 style={{ color: '#7af', margin: '0 0 0.75rem 0' }}>ΔV Metrics</h3>
      <table style={{ borderCollapse: 'collapse', width: '100%' }}>
        <tbody>
          <MetricRow label="Departure MJD" value={metrics.dep_mjd.toFixed(1)} />
          <MetricRow label="Arrival MJD" value={metrics.arr_mjd.toFixed(1)} />
          <MetricRow label="TOF" value={`${metrics.tof_days.toFixed(0)} days`} />
          <MetricRow label="ΔV Total" value={`${metrics.delta_v_total.toFixed(3)} km/s`} highlight />
        </tbody>
      </table>
    </div>
  );
};

function MetricRow({ label, value, highlight }: { label: string; value: string; highlight?: boolean }) {
  return (
    <tr>
      <td style={{ color: '#aaa', paddingBottom: '0.4rem', paddingRight: '0.5rem' }}>{label}</td>
      <td style={{ color: highlight ? '#ffdd44' : '#eee', fontWeight: highlight ? 700 : 400 }}>
        {value}
      </td>
    </tr>
  );
}

export default DeltaVMetrics;
