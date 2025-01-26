"use server";

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
    console.log("res: ", response);
    if (!response.ok) {
      const errorMsg = await response.text();
      return errorMsg;
    }
    const data = await response.json();
    return data;
  } catch (err) {
    console.log("requestPOST Error: ", err);
  }
};
