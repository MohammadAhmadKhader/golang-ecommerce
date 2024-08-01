CREATE Table IF NOT EXISTS products (
    `id` INT UNSIGNED AUTO_INCREMENT NOT NULL,
    `name` varchar(255) NOT NULL,
    `description` TEXT NOT NULL,
    `image` varchar(255) NOT NULL,
    `price` DECIMAL(10,2) NOT NULL,
    `quantity` INT UNSIGNED NOT NULL,
    `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updatedAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP On Update CURRENT_TIMESTAMP,

    PRIMARY KEY(`id`)
);