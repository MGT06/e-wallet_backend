INSERT INTO payment_methods (name, type, icon_url, tax_percent, created_at, xendit_code) VALUES
('BCA', 'bank', 'https://example.com/icons/bca.png', 10.00, NOW(), 'BCA_VIRTUAL_ACCOUNT'),
('BRI', 'bank', 'https://example.com/icons/mandiri.png', 10.50, NOW(), 'BRI_VIRTUAL_ACCOUNT'),
('GoPay', 'ewallet', 'https://example.com/icons/gopay.png', 3.00, NOW(), 'GOPAY'),
('OVO', 'ewallet', 'https://example.com/icons/ovo.png', 4.00, NOW(), 'OVO'),
('Dana', 'ewallet', 'https://example.com/icons/dana.png', 2.50, NOW(), 'DANA'); 