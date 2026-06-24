// ============================================
// Feature Plugin System — Registry & Types
// ============================================
// New features only need a config file + component
// in their own directory. No layout changes needed.
// ============================================

export type FeaturePosition =
  | "global-header"
  | "global-footer"
  | "global-floating"
  | "home-before-list"
  | "home-sidebar"
  | "home-after-list"
  | "post-before-content"
  | "post-sidebar"
  | "post-after-content"
  | "post-footer";

export interface BlogFeature {
  /** Unique feature key, e.g. "writing-heatmap" */
  key: string;
  /** Display name for admin panel */
  name: string;
  /** Whether this feature is active */
  enabled: boolean;
  /** Where the feature mounts */
  position: FeaturePosition;
  /** Component name (must match the .astro or .vue file in the feature dir) */
  component: string;
  /** Client loading strategy for interactive islands */
  client?: "load" | "idle" | "visible" | "media" | "only";
  /** Feature-specific configuration */
  config?: Record<string, unknown>;
}

// ============================================
// Feature Registry
// ============================================
// To add a new feature:
// 1. Create src/features/<name>/config.ts
// 2. Import and register below
// 3. Create your component
// ============================================

import { heatmapFeature } from "./heatmap/config";
import { petFeature } from "./pet/config";
import { commentFeature } from "./comment/config";
import { musicFeature } from "./music/config";
import { analyticsFeature } from "./analytics/config";

const features: BlogFeature[] = [
  heatmapFeature,
  petFeature,
  commentFeature,
  musicFeature,
  analyticsFeature,
];

/** Get all enabled features for a given slot position */
export function getEnabledFeatures(position: FeaturePosition): BlogFeature[] {
  return features.filter((f) => f.enabled && f.position === position);
}

/** Get a feature by key */
export function getFeature(key: string): BlogFeature | undefined {
  return features.find((f) => f.key === key);
}

/** All registered features (for admin UI) */
export function getAllFeatures(): BlogFeature[] {
  return features;
}
