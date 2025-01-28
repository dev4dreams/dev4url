"use server";
import * as Sentry from "@sentry/nextjs";

export const POST = async ({
  url,
  headers = { "Content-Type": "application/json" },
  body = null,
}) => {
  try {
    const response = await fetch(url, {
      method: "POST",
      headers,
      body: body ? JSON.stringify(body) : null,
    });

    if (!response.ok) {
      const errorMsg = await response.text();
      return errorMsg;
    }
    const data = await response.json();
    return data;
  } catch (err) {
    console.error("requestPOST Error: ", err);
    Sentry.captureException(err, "POST Error");
  }
};
