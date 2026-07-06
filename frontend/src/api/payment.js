import api from "./api";

export async function pay(
  orderId,
  idempotencyKey,
) {
  const response = await api.post(
    "/payments",
    {
      order_id: orderId,
    },
    {
      headers: {
        "Idempotency-Key": idempotencyKey,
      },
    },
  );

  return response.data;
}