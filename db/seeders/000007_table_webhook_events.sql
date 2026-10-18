INSERT INTO webhook_events (event_id, event_status, payload) VALUES
(
  'evt_xen_inv_created_001',
  'INVOICE.CREATED',
  '{"id": "inv_xen_9a8b7c6d5e4f", "status": "PENDING", "merchant_id": "merchant_123", "amount": 150000, "created": "2026-10-17T01:00:00Z"}'::jsonb
),
(
  'evt_xen_inv_paid_002',
  'INVOICE.PAID',
  '{"id": "inv_xen_1a2b3c4d5e6f", "status": "PAID", "merchant_id": "merchant_123", "amount": 250000, "paid_at": "2026-10-17T01:15:00Z", "payment_method": "QRIS"}'::jsonb
),
(
  'evt_xen_inv_expired_003',
  'INVOICE.EXPIRED',
  '{"id": "inv_xen_112233445566", "status": "EXPIRED", "merchant_id": "merchant_123", "amount": 50000, "expired_at": "2026-10-16T18:45:00Z"}'::jsonb
),
(
  'evt_xen_inv_created_004',
  'INVOICE.PENDING',
  '{"id": "inv_xen_77889900aabb", "status": "PENDING", "merchant_id": "merchant_123", "amount": 1000000, "created": "2026-10-17T02:00:00Z"}'::jsonb
),
(
  'evt_xen_inv_paid_005',
  'INVOICE.SETTLED',
  '{"id": "inv_xen_ffeeddccbbaa", "status": "SETTLED", "merchant_id": "merchant_123", "amount": 75000, "paid_at": "2026-10-17T02:10:00Z", "payment_method": "BCA"}'::jsonb
);