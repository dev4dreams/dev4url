export const POST = async ({
  url,
  headers = { "Content-Type": "application/json" },
  body = null,
}) => {
  try {
    console.log("requestPOST URL: ", url);
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
    console.log("POST DATA: ", data);
    return data;
  } catch (err) {
    console.log("requestPOST Error: ", err);
  }
};
