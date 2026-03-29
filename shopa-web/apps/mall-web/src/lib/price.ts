"use client";

function normalizeCents(value: number | string | bigint): bigint {
  if (typeof value === "bigint") {
    return value;
  }
  if (typeof value === "number") {
    if (!Number.isFinite(value)) {
      return 0n;
    }
    return BigInt(Math.trunc(value));
  }
  const trimmed = value.trim();
  if (!trimmed) {
    return 0n;
  }
  try {
    return BigInt(trimmed);
  } catch {
    return 0n;
  }
}

function withThousands(rawInteger: string): string {
  return rawInteger.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

export function formatCnyFromCents(value: number | string | bigint): string {
  const cents = normalizeCents(value);
  const negative = cents < 0n;
  const abs = negative ? -cents : cents;

  const integerPart = abs / 100n;
  const fractionalPart = abs % 100n;

  const integerText = withThousands(integerPart.toString());
  const fractionalText = fractionalPart.toString().padStart(2, "0");

  return `${negative ? "-" : ""}${integerText}.${fractionalText}`;
}

