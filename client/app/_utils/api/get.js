"use server";

import * as Sentry from "@sentry/nextjs";

export const GET = async (url, params = {}) => {
  try {
    // Add query parameters if they exist
    const queryString = new URLSearchParams(params).toString();
    const fullUrl = queryString ? `${url}?${queryString}` : url;

    const response = await fetch(fullUrl, {
      method: "GET",
      headers: {
        Accept: "application/json",
      },
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText);
    }

    return await response.json();
  } catch (error) {
    Sentry.captureException(err, "POST Error");

    throw error;
  }
};
