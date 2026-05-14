import React from 'react';

interface Props {
  depStart: number; depEnd: number;
  arrStart: number; arrEnd: number;
  onDepStartChange: (v: number) => void;
  onDepEndChange: (v: number) => void;
  onArrStartChange: (v: number) => void;
  onArrEndChange: (v: number) => void;
  onCompute: () => void;
}

const DateSelector: React.FC<Props> = ({
  depStart, depEnd, arrStart, arrEnd,
  onDepStartChange, onDepEndChange,
  onArrStartChange, onArrEndChange,
  onCompute,
}) => {
  return (
    <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap', marginBottom: '1rem', alignItems: 'flex-end' }}>
      <MJDField label="Departure Start (MJD)" value={depStart} onChange={onDepStartChange} />
      <MJDField label="Departure End (MJD)" value={depEnd} onChange={onDepEndChange} />
      <MJDField label="Arrival Start (MJD)" value={arrStart} onChange={onArrStartChange} />
      <MJDField label="Arrival End (MJD)" value={arrEnd} onChange={onArrEndChange} />
      <button
        onClick={onCompute}
        style={{
          background: '#334', border: '1px solid #7af', color: '#7af',
          padding: '0.5rem 1.5rem', borderRadius: 4, cursor: 'pointer', fontFamily: 'monospace',
        }}
      >
        Compute Grid
      </button>
    </div>
  );
};

function MJDField({ label, value, onChange }: { label: string; value: number; onChange: (v: number) => void }) {
  return (
    <label style={{ display: 'flex', flexDirection: 'column', gap: 4, color: '#aaa', fontSize: '0.8rem' }}>
      {label}
      <input
        type="number"
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
        style={{ background: '#1a1a2e', border: '1px solid #334', color: '#eee', padding: '0.3rem', borderRadius: 4, width: 120 }}
      />
    </label>
  );
}

export default DateSelector;
