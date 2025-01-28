export async function createShortUrl({ originalUrl, customUrl, user }) {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_TEST_URL}/createUrl`,
      {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          originalUrl,
          customUrl,
          user,
          createTime: new Date().toISOString(),
        }),
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error("Failed to create short URL:", error);
    throw error; // Re-throw to handle in the component
  }
}
