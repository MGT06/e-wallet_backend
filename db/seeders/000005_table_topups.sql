INSERT INTO topups (reference_no, account_id, payment_method_id, amount, tax, total_amount, status, paid_at, created_at) VALUES
('TOPUP-20260601-001', 1, 1, 1000000, 1000, 1001000, 'success', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
('TOPUP-20260602-002', 2, 3, 500000, 1000, 508500, 'success', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
('TOPUP-20260603-003', 3, 4, 300000, 1000, 305500, 'success', NOW() - INTERVAL '5 hours', NOW() - INTERVAL '5 hours'),
('TOPUP-20260603-004', 4, 2, 2000000, 2000, 2002000, 'pending', NULL, NOW());