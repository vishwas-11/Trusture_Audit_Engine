// src/components/BlockchainScene.tsx
import { useRef } from "react";
import { Canvas, useFrame } from "@react-three/fiber";
import { Stars, Float } from "@react-three/drei";
import * as THREE from "three";

const FloatingBlock = () => {
  const meshRef = useRef<THREE.Mesh>(null);

  // Rotate the mesh every frame
  useFrame((_, delta) => {
    if (meshRef.current) {
      meshRef.current.rotation.x -= delta * 0.2;
      meshRef.current.rotation.y -= delta * 0.2;
    }
  });

  return (
    <Float speed={2} rotationIntensity={0.5} floatIntensity={1}>
      <mesh ref={meshRef} scale={2.5}>
        <icosahedronGeometry args={[1, 1]} />
        {/* Wireframe gives that "tech/blockchain" structure look */}
        <meshStandardMaterial
          color="#38bdf8" // Tailwind cyan-400
          wireframe
          emissive="#0ea5e9"
          emissiveIntensity={0.5}
        />
      </mesh>
    </Float>
  );
};

const BlockchainScene = () => {
  return (
    <div className="absolute inset-0 -z-10 bg-slate-950">
      <Canvas camera={{ position: [0, 0, 6] }}>
        {/* Lighting */}
        <ambientLight intensity={0.5} />
        <pointLight position={[10, 10, 10]} intensity={1} color="#38bdf8" />
        
        {/* Background Stars for depth */}
        <Stars radius={100} depth={50} count={5000} factor={4} saturation={0} fade speed={1} />
        
        {/* The Main 3D Object */}
        <FloatingBlock />
      </Canvas>
    </div>
  );
};

export default BlockchainScene;