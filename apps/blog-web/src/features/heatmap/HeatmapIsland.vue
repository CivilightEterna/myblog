<template>
  <div class="heatmap-grid" ref="gridRef">
    <div class="heatmap-months">
      <span v-for="(m, i) in monthLabels" :key="i" class="month-label">{{ m }}</span>
    </div>
    <div class="heatmap-body">
      <div class="heatmap-weekdays">
        <span class="wd-label">一</span>
        <span class="wd-label">三</span>
        <span class="wd-label">五</span>
      </div>
      <div class="heatmap-cells">
        <div
          v-for="(cell, idx) in cells"
          :key="idx"
          class="heatmap-cell"
          :class="`level-${cell.level}`"
          :title="cell.tooltip"
          @mouseenter="showTooltip($event, cell)"
          @mouseleave="hideTooltip"
        >
          <span class="sr-only">{{ cell.tooltip }}</span>
        </div>
      </div>
    </div>
    <!-- Tooltip -->
    <Teleport to="body">
      <div
        v-if="tooltip.visible"
        class="heatmap-tooltip"
        :style="{ left: tooltip.x + 'px', top: tooltip.y + 'px' }"
      >
        <div class="tooltip-date">{{ tooltip.date }}</div>
        <div class="tooltip-count">{{ tooltip.count }} 篇文章</div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";

interface HeatmapCell {
  date: string;
  count: number;
  level: number;
  tooltip: string;
  dayOfWeek: number;
}

const cells = ref<HeatmapCell[]>([]);
const monthLabels = ref<string[]>([]);
const tooltip = ref({ visible: false, x: 0, y: 0, date: "", count: 0 });

onMounted(async () => {
  try {
    const res = await fetch("/data/writing-heatmap.json");
    const data: { date: string; count: number }[] = await res.json();

    // Build a map of date -> count
    const countMap = new Map<string, number>();
    for (const d of data) {
      countMap.set(d.date, d.count);
    }

    // Find max count for color scaling
    const maxCount = Math.max(1, ...data.map((d) => d.count));

    // Generate cells for the last 365 days
    const today = new Date();
    const result: HeatmapCell[] = [];
    const months: string[] = [];

    for (let i = 364; i >= 0; i--) {
      const d = new Date(today);
      d.setDate(d.getDate() - i);
      const dateStr = d.toISOString().slice(0, 10);
      const count = countMap.get(dateStr) || 0;
      const level = count === 0 ? 0 : Math.min(4, Math.ceil((count / maxCount) * 4));

      result.push({
        date: dateStr,
        count,
        level,
        tooltip: `${dateStr}\n${count} 篇`,
        dayOfWeek: d.getDay(),
      });

      // Track month boundaries
      if (i === 364 || d.getDate() === 1) {
        months.push(`${d.getMonth() + 1}月`);
      } else if (d.getDate() <= 3 && i < 360) {
        months.push("");
      }
    }

    cells.value = result;
    monthLabels.value = months;
  } catch (e) {
    console.error("Failed to load heatmap data:", e);
  }
});

function showTooltip(event: MouseEvent, cell: HeatmapCell) {
  tooltip.value = {
    visible: true,
    x: event.clientX + 12,
    y: event.clientY - 40,
    date: cell.date,
    count: cell.count,
  };
}

function hideTooltip() {
  tooltip.value.visible = false;
}
</script>
