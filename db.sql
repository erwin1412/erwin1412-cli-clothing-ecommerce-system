CREATE DATABASE pairing;

use pairing;

CREATE TABLE users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL
);

CREATE TABLE user_profiles (
  id INT PRIMARY KEY AUTO_INCREMENT,
  userId INT NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  FOREIGN KEY (userId) REFERENCES users(id)
);

CREATE TABLE cloths (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  price INT NOT NULL,
  size VARCHAR(50) NOT NULL,
  color VARCHAR(50) NOT NULL,
  stock INT NOT NULL,
  description TEXT
);

CREATE TABLE transactions (
  id INT PRIMARY KEY AUTO_INCREMENT,
  userId INT NOT NULL,
  total INT NOT NULL,
  FOREIGN KEY (userId) REFERENCES users(id)
);

CREATE TABLE transaction_details (
  id INT PRIMARY KEY AUTO_INCREMENT,
  transactionId INT NOT NULL,
  clothId INT NOT NULL,
  price INT NOT NULL,
  qty INT NOT NULL,
  total INT NOT NULL,
  FOREIGN KEY (transactionId) REFERENCES transactions(id),
  FOREIGN KEY (clothId) REFERENCES cloths(id)
);

CREATE TABLE user_cloth_favorites (
  userId INT NOT NULL,
  clothId INT NOT NULL,
  PRIMARY KEY (userId, clothId),
  FOREIGN KEY (userId) REFERENCES users(id),
  FOREIGN KEY (clothId) REFERENCES cloths(id)
);



INSERT INTO cloths (name, price, size, color, stock, description) VALUES
('T-Shirt Basic', 100000, 'M', 'White', 50, 'Basic white cotton t-shirt'),
('Denim Jacket', 250000, 'L', 'Blue', 30, 'Classic blue denim jacket'),
('Hoodie Oversize', 200000, 'XL', 'Black', 20, 'Oversized black hoodie with pocket'),
('Chino Pants', 180000, 'L', 'Khaki', 40, 'Comfortable khaki chino pants'),
('Flannel Shirt', 150000, 'M', 'Red', 25, 'Red and black flannel shirt'),
('Polo Shirt', 130000, 'S', 'Navy', 35, 'Casual navy polo shirt'),
('Jeans Skinny Fit', 220000, 'M', 'Dark Blue', 28, 'Skinny fit dark blue jeans'),
('Sweater Knit', 175000, 'L', 'Gray', 18, 'Warm gray knit sweater'),
('Shorts Cargo', 120000, 'M', 'Green', 22, 'Green cargo shorts with pockets'),
('Long Sleeve Tee', 110000, 'L', 'Maroon', 27, 'Maroon long sleeve cotton tee');
