INSERT INTO accounts (fullname, image_url, phone, email, password, pin, balance, email_verified_at, created_at, updated_at) VALUES
('Budi Santoso', 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde', '081234567890', 'budi.santoso@example.com', '$2a$12$LQv3c1yqSNWvklHnWTItze5zQ2K1vO0x6Vn7PjG2h1qR8yF3V6x2W', '$2a$12$LQv3c1yqSNWvklHnWTItze5zQ2K1vO0x6Vn7PjG2h1qR8yF3V6x2W', 1500000, NOW(), NOW(), NOW()),
('Siti Aminah', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330', '081298765432', 'siti.aminah@example.com', '$2a$12$LQv3c1yqSNWvklHnWTItze5zQ2K1vO0x6Vn7PjG2h1qR8yF3V6x2W', '$2a$12$LQv3c1yqSNWvklHnWTItze5zQ2K1vO0x6Vn7PjG2h1qR8yF3V6x2W', 750000, NOW(), NOW(), NOW()),
('Ahmad Fauzi', 'https://images.unsplash.com/photo-1570295999919-56ceb5ecca61', '081388889999', 'ahmad.fauzi@example.com', '$2a$12$LQv3c1yqSNWvklHnWTItze5zQ2K1vO0x6Vn7PjG2h1qR8yF3V6x2W', '$2a$12$LQv3c1yqSNWvklHnWTItze5zQ2K1vO0x6Vn7PjG2h1qR8yF3V6x2W', 300000, NOW(), NOW(), NOW()),
('Dewi Lestari', 'https://images.unsplash.com/photo-1580489944761-15a19d654956', '081577776666', 'dewi.lestari@example.com', '$2a$12$LQv3c1yqSNWvklHnWTItze5zQ2K1vO0x6Vn7PjG2h1qR8yF3V6x2W', '$2a$12$LQv3c1yqSNWvklHnWTItze5zQ2K1vO0x6Vn7PjG2h1qR8yF3V6x2W', 2500000, NOW(), NOW(), NOW());

INSERT INTO payment_methods (name, type, icon_url, tax_percent, created_at) VALUES
('BCA', 'bank', 'https://example.com/icons/bca.png', 10.00, NOW()),
('Mandiri', 'bank', 'https://example.com/icons/mandiri.png', 10.50, NOW()),
('GoPay', 'ewallet', 'https://example.com/icons/gopay.png', 3.00, NOW()),
('OVO', 'ewallet', 'https://example.com/icons/ovo.png', 4.00, NOW()),
('Dana', 'ewallet', 'https://example.com/icons/dana.png', 2.50, NOW());

INSERT INTO transfers (reference_no, sender_id, receiver_id, amount, notes, status, created_at) VALUES
('TRX-20260601-001', 1, 2, 100000, 'Bayar patungan makanan', 'success', NOW() - INTERVAL '1 day'),
('TRX-20260602-002', 2, 3, 50000, 'Ganti uang kopi', 'success', NOW() - INTERVAL '12 hours'),
('TRX-20260603-003', 4, 1, 250000, 'Cicilan buku', 'success', NOW() - INTERVAL '2 hours');

INSERT INTO topups (reference_no, account_id, payment_method_id, amount, tax, total_amount, status, paid_at, created_at) VALUES
('TOPUP-20260601-001', 1, 1, 1000000, 1000, 1001000, 'success', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
('TOPUP-20260602-002', 2, 3, 500000, 1000, 508500, 'success', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
('TOPUP-20260603-003', 3, 4, 300000, 1000, 305500, 'success', NOW() - INTERVAL '5 hours', NOW() - INTERVAL '5 hours'),
('TOPUP-20260603-004', 4, 2, 2000000, 2000, 2002000, 'pending', NULL, NOW());

INSERT INTO transaction (account_id, type, amount, balance_before, balance_after, topup_id, transfer_id, created_at) VALUES
(1, 'topup', 1000000, 500000, 1500000, 1, NULL, NOW() - INTERVAL '2 days'),
(1, 'transfer', -100000, 1600000, 1500000, NULL, 1, NOW() - INTERVAL '1 day'),
(2, 'topup', 500000, 250000, 750000, 2, NULL, NOW() - INTERVAL '1 day'),
(2, 'transfer', -50000, 800000, 750000, NULL, 2, NOW() - INTERVAL '12 hours'),
(3, 'topup', 300000, 0, 300000, 3, NULL, NOW() - INTERVAL '5 hours'),
(4, 'transfer', -250000, 2750000, 2500000, NULL, 3, NOW() - INTERVAL '2 hours');

INSERT INTO favorites (account_id, favorite_account_id, created_at) VALUES
(1, 2, NOW()),
(1, 3, NOW()),
(2, 1, NOW()),
(4, 1, NOW());