/**
 * Calculate estimated reading time from word count or raw text.
 * Assumes ~300 Chinese characters or ~200 English words per minute.
 */

const CHARS_PER_MINUTE = 300;
const WORDS_PER_MINUTE = 200;

/**
 * Count Chinese characters in a string.
 */
function countChineseChars(text: string): number {
  const matches = text.match(/[一-鿿㐀-䶿]/g);
  return matches ? matches.length : 0;
}

/**
 * Count English words in a string.
 */
function countEnglishWords(text: string): number {
  // Remove Chinese characters, then split on whitespace
  const englishOnly = text.replace(/[一-鿿㐀-䶿]/g, "");
  const words = englishOnly.trim().split(/\s+/);
  return words[0] === "" ? 0 : words.length;
}

/**
 * Estimate reading time in minutes.
 * Returns at least 1 minute.
 */
export function estimateReadingTime(content: string): number {
  const chineseChars = countChineseChars(content);
  const englishWords = countEnglishWords(content);

  const minutes =
    chineseChars / CHARS_PER_MINUTE + englishWords / WORDS_PER_MINUTE;

  return Math.max(1, Math.ceil(minutes));
}

/**
 * Format reading time display.
 */
export function formatReadingTime(content: string): string {
  const minutes = estimateReadingTime(content);
  if (minutes < 60) {
    return `${minutes} 分钟`;
  }
  const hours = Math.floor(minutes / 60);
  const remaining = minutes % 60;
  return remaining > 0 ? `${hours} 小时 ${remaining} 分钟` : `${hours} 小时`;
}
