import { createShortUrl } from "../request/create";
import { getAllUrls } from "../request/getAllUrl";

export const url = {
  GET: {
    url: "",
    all: getAllUrls,
  },
  POST: {
    create: createShortUrl,
  },
};
