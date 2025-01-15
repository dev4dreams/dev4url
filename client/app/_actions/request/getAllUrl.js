export async function getAllUrls() {
  try {
    const res = await fetch(`${process.env.NEXT_PUBLIC_TEST_URL}/urls`, {
      method: "GET",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    });

    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`);
    }

    return await res.json();
  } catch (error) {
    console.error("Failed to fetch URLs:", error);
    throw error;
  }
}
