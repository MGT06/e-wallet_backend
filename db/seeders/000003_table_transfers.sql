INSERT INTO transfers (reference_no, sender_id, receiver_id, amount, notes, status, created_at) VALUES
('TRX-20260601-001', 1, 2, 100000, 'Bayar patungan makanan', 'success', NOW() - INTERVAL '1 day'),
('TRX-20260602-002', 2, 3, 50000, 'Ganti uang kopi', 'success', NOW() - INTERVAL '12 hours'),
('TRX-20260603-003', 4, 1, 250000, 'Cicilan buku', 'success', NOW() - INTERVAL '2 hours');