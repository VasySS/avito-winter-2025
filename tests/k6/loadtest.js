import http from "k6/http";
import { check } from "k6";
import { SharedArray } from "k6/data";

// ignore 4xx
http.setResponseCallback(
  http.expectedStatuses(400, 404, { min: 200, max: 299 })
);

const users = new SharedArray("users", function () {
  return JSON.parse(open("./users.json"));
});

const userTokens = new Map();

const items = ["t-shirt", "pen", "cup", "socks"];

const baseURL = "http://localhost:8080";

export const options = {
  stages: [
    { duration: "10s", target: 50 },
    { duration: "10s", target: 100 },
    { duration: "30", target: 150 },
    { duration: "10s", target: 0 },
  ],
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<500"],
  },
};

export default function () {
  const user = users[Math.floor(Math.random() * users.length)];
  let token = userTokens.get(user.username);

  if (!token) {
    const authRes = http.post(
      `${baseURL}/api/auth`,
      JSON.stringify({
        username: user.username,
        password: user.password,
      }),
      { headers: { "Content-Type": "application/json" } }
    );

    if (
      !check(authRes, {
        "auth status 200": (r) => r.status === 200,
        "auth token received": (r) => r.json().token !== undefined,
      })
    ) {
      return;
    }

    token = authRes.json().token;
    userTokens.set(user.username, token);
  }

  const headers = {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  };

  const infoRes = http.get(`${baseURL}/api/info`, headers);
  check(infoRes, { "info status 200": (r) => r.status === 200 });

  let userReciever;
  do {
    userReciever = users[Math.floor(Math.random() * users.length)];
  } while (userReciever.username === user.username);

  const sendCoinRes = http.post(
    `${baseURL}/api/send-coin`,
    JSON.stringify({
      toUser: userReciever.username,
      amount: 10,
    }),
    headers
  );

  //   console.log(sendCoinRes);

  check(sendCoinRes, { "sendCoin status != 500": (r) => r.status !== 500 });

  const item = items[Math.floor(Math.random() * items.length)];
  const buyRes = http.post(`${baseURL}/api/buy/${item}`, null, headers);

  //   console.log(buyRes);

  check(buyRes, { "buyItem status != 500": (r) => r.status !== 500 });
}
