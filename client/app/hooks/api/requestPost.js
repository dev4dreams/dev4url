export const requestPost = async ({
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
    console.log("TEST: ", response);
    return await response.json();
  } catch (err) {
    console.log("requestPOST Error: ", err);
  }
};
