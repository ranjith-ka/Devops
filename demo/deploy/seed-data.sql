INSERT INTO cases (title, details, sender, eta, sla_days, hypercare, label, hs_code, preference, supp_units, assigned_to, priority_score, created_at) VALUES
-- High priority: premium customer, ETA soon, hypercare
('Container delayed at customs', 'Container XYZ stuck at port; customer requests urgent clearance', 'ops@premium-customer.com', '2025-09-13 10:00:00+00', 1, true, 'customs_issue', 'UNKNOWN', 'Preferential', 'EA', 'alice', 95, now()),
('Perishable goods held - urgent', 'Refrigerated goods held — risk of spoilage. Immediate action required', 'logistics@vipclient.com', '2025-09-13 08:00:00+00', 0, true, 'perishable_hold', '1602.49', 'TemperatureSensitive', 'KG', 'bob', 98, now()),

-- Medium priority: short SLA but not hypercare
('Documentation missing for shipment ABC', 'Missing commercial invoice and packing list for AWB123', 'customer@example.com', '2025-09-16 12:00:00+00', 2, false, 'missing_docs', 'UNKNOWN', '', '', 'charlie', 60, now()),
('Customs classification query', 'Customer asks clarification on HS for electric battery units', 'trader@electronics.com', '2025-09-18 00:00:00+00', 3, false, 'hs_query', '8507.80', '', 'EA', NULL, 50, now()),

-- Low priority: longer ETA, non-hypercare
('Routine shipment update', 'Standard shipment update request; no urgent SLAs', 'user@smallclient.com', '2025-09-30 09:00:00+00', 7, false, 'general', 'UNKNOWN', '', '', NULL, 10, now()),
('Bulk items for catalog', 'Clothing consignments for seasonal catalog; multiple items', 'buyer@retailco.com', '2025-10-03 15:00:00+00', 10, false, 'catalog_shipping', '6203.42', '', 'PCS', NULL, 8, now()),

-- Cases to demonstrate duplicate detection (similar titles / same sender)
('Shipment delayed at customs - Container XYZ', 'Follow-up: same container reported delayed', 'ops@premium-customer.com', '2025-09-13 11:00:00+00', 1, true, 'customs_issue', 'UNKNOWN', 'Preferential', 'EA', 'dave', 92, now()),
('Container XYZ - customs hold', 'New email thread referencing Container XYZ; likely duplicate', 'ops@premium-customer.com', '2025-09-13 11:30:00+00', 1, true, 'customs_issue', 'UNKNOWN', 'Preferential', 'EA', NULL, 90, now()),

-- Examples for bulk update and HS suggestions
('Electronics shipment - batteries included', 'Contains rechargeable batteries and chargers', 'supplier@electronix.com', '2025-09-20 09:00:00+00', 5, false, 'electronics', '8507.80', '', 'EA', NULL, 40, now()),
('Assorted toys shipment', 'Multiple toy items for retail; includes small electronic parts', 'distributor@toys.com', '2025-09-22 10:00:00+00', 4, false, 'toys', '9503.00', '', 'PCS', NULL, 30, now());
