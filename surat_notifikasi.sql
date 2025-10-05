-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1:3308
-- Generation Time: Oct 05, 2025 at 12:35 PM
-- Server version: 10.4.32-MariaDB
-- PHP Version: 8.2.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `surat_notifikasi`
--

-- --------------------------------------------------------

--
-- Table structure for table `letters`
--

CREATE TABLE `letters` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `type_id` int(11) NOT NULL,
  `status` enum('pending','accepted','rejected') DEFAULT 'pending',
  `reject_reason` text DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `letters`
--

INSERT INTO `letters` (`id`, `user_id`, `type_id`, `status`, `reject_reason`, `created_at`, `updated_at`) VALUES
(5, 5, 1, 'accepted', '', '2025-10-03 07:07:29', '2025-10-04 23:04:16'),
(7, 5, 1, 'pending', '', '2025-10-04 22:55:46', '2025-10-04 22:55:46'),
(8, 5, 2, 'rejected', 'Ditolak utk testing 5 Okt 2025', '2025-10-04 22:59:11', '2025-10-04 23:06:40');

-- --------------------------------------------------------

--
-- Table structure for table `letter_types`
--

CREATE TABLE `letter_types` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `description` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `letter_types`
--

INSERT INTO `letter_types` (`id`, `name`, `description`) VALUES
(1, 'Surat A', 'Deskripsi Surat A'),
(2, 'Surat B', 'Deskripsi Surat B'),
(3, 'Surat C', 'Deskripsi Surat C');

-- --------------------------------------------------------

--
-- Table structure for table `roles`
--

CREATE TABLE `roles` (
  `id` int(11) NOT NULL,
  `name` varchar(191) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `roles`
--

INSERT INTO `roles` (`id`, `name`) VALUES
(1, 'admin'),
(2, 'reviewer'),
(3, 'user');

-- --------------------------------------------------------

--
-- Table structure for table `settings`
--

CREATE TABLE `settings` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `telegram_chatid` varchar(50) DEFAULT NULL,
  `wa_number` varchar(20) DEFAULT NULL,
  `allow_telegram` enum('yes','no') DEFAULT 'no',
  `allow_wa` enum('yes','no') DEFAULT 'no',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `settings`
--

INSERT INTO `settings` (`id`, `user_id`, `telegram_chatid`, `wa_number`, `allow_telegram`, `allow_wa`, `created_at`, `updated_at`) VALUES
(5, 5, '123456789', '628123456789', 'yes', 'yes', '2025-10-04 21:44:03', '2025-10-05 05:51:15'),
(6, 8, '123456789', '628123456789', 'yes', 'yes', '2025-10-04 21:44:17', '2025-10-05 05:51:11');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `role_id`, `name`, `email`, `password`, `created_at`, `updated_at`) VALUES
(4, 1, 'Admin Satu', 'admin@mail.com', '$2a$10$0rI.bujqzSiOp.ctWcDxsunqN7NpHTCOVpkM/zoEGm35HX5a0qIMy', '2025-10-03 03:21:12', '2025-10-03 11:25:55'),
(5, 3, 'User Baru', 'user@mail.com', '$2a$10$rpRTdYg3x.Y2N.GJynKMrOkFDYKPdBP5OLvQo8UTPkE2h4JlboC.m', '2025-10-03 03:45:41', '2025-10-04 21:55:31'),
(7, 3, 'Admin Dua', 'admin2@mail.com', '$2a$10$.n7BlFnD.VIPKXBH9cIWjuWPcpg32bIptO/8GXhqsVq24mcWncpSK', '2025-10-03 09:30:26', '2025-10-03 09:30:26'),
(8, 2, 'Reviewer', 'reviewer@mail.com', '$2a$10$ML55BitZFoysUl057BOXz.q/tmPMUs96nUThWGgr5883DkuI5aY7e', '2025-10-04 08:33:48', '2025-10-04 21:56:31');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `letters`
--
ALTER TABLE `letters`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`),
  ADD KEY `type_id` (`type_id`);

--
-- Indexes for table `letter_types`
--
ALTER TABLE `letter_types`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `roles`
--
ALTER TABLE `roles`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uni_roles_name` (`name`);

--
-- Indexes for table `settings`
--
ALTER TABLE `settings`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD KEY `role_id` (`role_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `letters`
--
ALTER TABLE `letters`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `letter_types`
--
ALTER TABLE `letter_types`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `roles`
--
ALTER TABLE `roles`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `settings`
--
ALTER TABLE `settings`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `letters`
--
ALTER TABLE `letters`
  ADD CONSTRAINT `letters_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `letters_ibfk_2` FOREIGN KEY (`type_id`) REFERENCES `letter_types` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `settings`
--
ALTER TABLE `settings`
  ADD CONSTRAINT `settings_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `users`
--
ALTER TABLE `users`
  ADD CONSTRAINT `users_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
