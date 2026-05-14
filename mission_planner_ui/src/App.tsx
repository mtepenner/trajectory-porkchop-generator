import React, { useState } from 'react';
import DateSelector from './components/DateSelector';
import PorkchopPlot from './components/PorkchopPlot';
import SolarSystem3D from './components/SolarSystem3D';
import DeltaVMetrics from './components/DeltaVMetrics';
import { useTrajectoryData } from './hooks/useTrajectoryData';

const App: React.FC = () => {
  const [depStart, setDepStart] = useState<number>(60000);
  const [depEnd, setDepEnd] = useState<number>(60200);
  const [arrStart, setArrStart] = useState<number>(60300);
  const [arrEnd, setArrEnd] = useState<number>(60700);

  const { gridData, trajectoryPath, metrics, loading, error, fetchGrid, fetchPath } =
    useTrajectoryData();

  const handleCompute = () => {
    fetchGrid(depStart, depEnd, arrStart, arrEnd, 100);
  };

  return (
    <div style={{ fontFamily: 'monospace', background: '#0a0a1a', color: '#e0e0ff', minHeight: '100vh', padding: '1rem' }}>
      <h1 style={{ color: '#7af', textAlign: 'center' }}>
        Interplanetary Trajectory Porkchop Generator
      </h1>

      <DateSelector
        depStart={depStart} depEnd={depEnd}
        arrStart={arrStart} arrEnd={arrEnd}
        onDepStartChange={setDepStart} onDepEndChange={setDepEnd}
        onArrStartChange={setArrStart} onArrEndChange={setArrEnd}
        onCompute={handleCompute}
      />

      {error && <p style={{ color: 'red' }}>{error}</p>}
      {loading && <p>Computing trajectories…</p>}

      <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
        {gridData && (
          <>
            <PorkchopPlot
              grid={gridData.grid}
              depAxis={gridData.dep_mjd_axis}
              arrAxis={gridData.arr_mjd_axis}
              onCellClick={(dep, arr) => fetchPath(dep, arr)}
            />
            <DeltaVMetrics metrics={metrics} />
          </>
        )}
      </div>

      {trajectoryPath && <SolarSystem3D path={trajectoryPath} />}
    </div>
  );
};

export default App;
