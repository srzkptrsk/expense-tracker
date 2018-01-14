# expense-tracker
Expense Tracker - golang powered telegram bot for simple expense tracker

Create tables in your MySQL database:

    CREATE TABLE `money` (`money_id` int(10) unsigned NOT NULL AUTO_INCREMENT, `user_id` int(10) DEFAULT NULL, `amount` decimal(10,2) NOT NULL, `category` varchar(255) CHARACTER SET latin1 DEFAULT NULL, `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`money_id`), UNIQUE KEY `money_id_UNIQUE` (`money_id`)) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;
    CREATE TABLE `user` (`user_id` int(10) unsigned NOT NULL, PRIMARY KEY (`user_id`), UNIQUE KEY `user_id_UNIQUE` (`user_id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;

For daemonize need make next steps:
* nano /lib/systemd/system/expensetracker.service
* *insert from expensetracker.service.sample with your paths*
* systemctl enable expensetracker.service
* service expensetracker start