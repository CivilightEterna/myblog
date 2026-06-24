/**
 * Format a date string (YYYY-MM-DD) to a human-readable format in Chinese.
 */
export function formatDate(dateStr: string, locale: string = "zh-CN"): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString(locale, {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

/**
 * Format a date string to ISO format.
 */
export function toISODate(dateStr: string): string {
  return new Date(dateStr).toISOString();
}

/**
 * Get the year from a date string.
 */
export function getYear(dateStr: string): number {
  return new Date(dateStr).getFullYear();
}

/**
 * Get the month (1-12) from a date string.
 */
export function getMonth(dateStr: string): number {
  return new Date(dateStr).getMonth() + 1;
}

/**
 * Group dates by year and month. Returns a Map of year -> month -> posts.
 */
export function groupByYearMonth<T>(
  items: T[],
  dateGetter: (item: T) => string
): Map<number, Map<number, T[]>> {
  const grouped = new Map<number, Map<number, T[]>>();

  for (const item of items) {
    const date = dateGetter(item);
    const year = getYear(date);
    const month = getMonth(date);

    if (!grouped.has(year)) {
      grouped.set(year, new Map());
    }
    const yearGroup = grouped.get(year)!;
    if (!yearGroup.has(month)) {
      yearGroup.set(month, []);
    }
    yearGroup.get(month)!.push(item);
  }

  return grouped;
}
