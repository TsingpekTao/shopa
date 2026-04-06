import { NextRequest, NextResponse } from "next/server";

type AssetReadUrlPayload = {
  url?: string;
};

type GoFrameResponse<T> = {
  code?: number;
  message?: string;
  data?: T;
};

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

function unwrapGoFrame<T>(payload: unknown): T {
  if (!payload || typeof payload !== "object") {
    return payload as T;
  }

  const response = payload as GoFrameResponse<T>;
  if (typeof response.code === "number" && "data" in response) {
    if (response.code !== 0) {
      throw new Error(response.message || "upstream request failed");
    }
    return response.data as T;
  }
  return payload as T;
}

async function issueAssetReadUrl(baseUrl: string, assetId: string, ttlSeconds: number) {
  try {
    const response = await fetch(
      `${baseUrl}/v1/media/assets/${encodeURIComponent(assetId)}/read-url?ttl_seconds=${ttlSeconds}`,
      {
        method: "GET",
        cache: "no-store"
      }
    );
    if (!response.ok) {
      return "";
    }

    const payload = unwrapGoFrame<AssetReadUrlPayload>(await response.json());
    return toString(payload?.url);
  } catch {
    return "";
  }
}

export async function GET(request: NextRequest) {
  const assetIds = Array.from(
    new Set(
      (request.nextUrl.searchParams.get("assetIds") || "")
        .split(",")
        .map((item) => item.trim())
        .filter((item) => item.length > 0)
    )
  ).slice(0, 50);
  const ttlSeconds = Math.min(Math.max(Number(request.nextUrl.searchParams.get("ttlSeconds") || "900"), 60), 3600);

  if (assetIds.length === 0) {
    return NextResponse.json({ items: [] });
  }

  const mediaBaseUrl = process.env.MEDIA_SERVICE_BASE_URL || "http://127.0.0.1:8006";
  const items = await Promise.all(
    assetIds.map(async (assetId) => ({
      assetId,
      url: await issueAssetReadUrl(mediaBaseUrl, assetId, ttlSeconds)
    }))
  );

  return NextResponse.json({
    items: items.filter((item) => item.url)
  });
}
