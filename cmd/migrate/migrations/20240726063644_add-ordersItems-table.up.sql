CREATE TABLE IF NOT EXISTS orderItems (
    `id` INT NOT NULL AUTO_INCREMENT,
    `orderId` INT UNSIGNED NOT NULL,
    `productId` INT UNSIGNED NOT NULL,
    `quantity` INT UNSIGNED NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,

    PRIMARY KEY(`id`),
    FOREIGN KEY(`orderId`) REFERENCES orders(`id`),
    FOREIGN KEY(`productId`) REFERENCES products(`id`)
);