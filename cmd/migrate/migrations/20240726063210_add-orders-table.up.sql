CREATE TABLE IF NOT EXISTS orders (
    `id` INT UNSIGNED AUTO_INCREMENT NOT NULL,
    `userId` INT UNSIGNED NOT NULL,
    `total` DECIMAL(10,2) NOT NULL,
    `status` ENUM('pending', 'completed', 'cancelled') NOT NULL DEFAULT 'pending',
    `address` TEXT NOT NULL,
    `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updatedAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP On Update CURRENT_TIMESTAMP,

    PRIMARY KEY(`id`),
    FOREIGN KEY(`userId`) REFERENCES users(`id`)
);