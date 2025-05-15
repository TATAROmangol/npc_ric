INSERT INTO institutions (name, inn, columns) VALUES
('A', 111, ARRAY['1', '2', '3']),
('B', 222, ARRAY['1', '2', '3']),
('C', 333, ARRAY['1', '2', '3']);

INSERT INTO mentors (name) VALUES
('A'),
('B'),
('C');

INSERT INTO forms (info, institution_id) VALUES
(ARRAY['A', 'B', 'C'], 1), 
(ARRAY['D', 'E', 'F'], 1), 
(ARRAY['G', 'H', 'I'], 1), 
(ARRAY['A', 'B', 'C'], 1),
(ARRAY['A', 'B', 'C'], 2);