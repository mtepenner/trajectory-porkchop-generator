import React, { useRef, useMemo } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, Line } from '@react-three/drei';
import * as THREE from 'three';

interface TrajectoryPoint { t: number; x: number; y: number; z: number; }
interface Props { path: TrajectoryPoint[]; }

const AU = 149597870.7;
const SCALE = 1 / AU; // normalise to AU for display

function TransferArc({ path }: Props) {
  const points = useMemo(
    () => path.map((p) => new THREE.Vector3(p.x * SCALE, p.z * SCALE, -p.y * SCALE)),
    [path]
  );
  return <Line points={points} color="#ffcc00" lineWidth={2} />;
}

function Planet({ radius, color, label }: { radius: number; color: string; label: string }) {
  const meshRef = useRef<THREE.Mesh>(null!);
  useFrame((_, delta) => { meshRef.current.rotation.y += delta * 0.5; });
  return (
    <mesh ref={meshRef} position={[radius, 0, 0]}>
      <sphereGeometry args={[0.03, 16, 16]} />
      <meshStandardMaterial color={color} />
    </mesh>
  );
}

const SolarSystem3D: React.FC<Props> = ({ path }) => {
  return (
    <div style={{ width: '100%', height: 400, background: '#050510', borderRadius: 8 }}>
      <Canvas camera={{ position: [0, 2, 4], fov: 50 }}>
        <ambientLight intensity={0.4} />
        <pointLight position={[0, 0, 0]} intensity={2} color="#fff5d0" />
        {/* Sun */}
        <mesh position={[0, 0, 0]}>
          <sphereGeometry args={[0.08, 32, 32]} />
          <meshStandardMaterial color="#ffdd00" emissive="#ffaa00" emissiveIntensity={1} />
        </mesh>
        {/* Earth orbit ring */}
        <mesh rotation={[Math.PI / 2, 0, 0]}>
          <ringGeometry args={[0.998, 1.002, 128]} />
          <meshBasicMaterial color="#4488ff" side={THREE.DoubleSide} transparent opacity={0.3} />
        </mesh>
        {/* Mars orbit ring */}
        <mesh rotation={[Math.PI / 2, 0, 0]}>
          <ringGeometry args={[1.522, 1.526, 128]} />
          <meshBasicMaterial color="#cc4422" side={THREE.DoubleSide} transparent opacity={0.3} />
        </mesh>
        <Planet radius={1.0} color="#3399ff" label="Earth" />
        <Planet radius={1.524} color="#cc5533" label="Mars" />
        <TransferArc path={path} />
        <OrbitControls />
      </Canvas>
    </div>
  );
};

export default SolarSystem3D;
