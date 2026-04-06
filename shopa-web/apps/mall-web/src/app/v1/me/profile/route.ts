import { NextRequest, NextResponse } from "next/server";

function resolveUpstreamBaseURL() {
  return process.env.NEXT_PUBLIC_API_BASE_URL?.trim() || process.env.MALL_API_BASE_URL?.trim() || "http://127.0.0.1:8000";
}

function buildForwardHeaders(request: NextRequest, hasBody: boolean): HeadersInit {
  const headers = new Headers();
  const authorization = request.headers.get("authorization");
  if (authorization) {
    headers.set("authorization", authorization);
  }
  if (hasBody) {
    headers.set("content-type", request.headers.get("content-type") || "application/json");
  }
  return headers;
}

async function proxyProfile(request: NextRequest, method: "GET" | "PATCH") {
  const upstreamUrl = new URL(`${resolveUpstreamBaseURL()}/v1/me/profile`);
  upstreamUrl.search = request.nextUrl.search;

  const hasBody = method !== "GET";
  const response = await fetch(upstreamUrl, {
    method,
    headers: buildForwardHeaders(request, hasBody),
    body: hasBody ? await request.text() : undefined,
    cache: "no-store"
  });

  const contentType = response.headers.get("content-type") || "application/json; charset=utf-8";
  return new NextResponse(await response.text(), {
    status: response.status,
    headers: {
      "content-type": contentType
    }
  });
}

export async function GET(request: NextRequest) {
  return proxyProfile(request, "GET");
}

export async function PATCH(request: NextRequest) {
  return proxyProfile(request, "PATCH");
}
