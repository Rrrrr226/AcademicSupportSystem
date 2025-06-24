DROP TABLE IF EXISTS `student_history_info`;
CREATE TABLE `student_history_info` (
    `school_year` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '学年',
    `semester` varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '学期',
    `student_id` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '学号',
    `student_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '学生姓名',
    `unit_id` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '学院ID',
    `unit_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '学院名称',
    `major_code` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '专业代码',
    `major` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '专业名称',
    `class_no` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '班级',
    `grade` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '年级',
    `teacher_id` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '辅导员',
    PRIMARY KEY (`school_year`,`semester`,`student_id`) USING BTREE,
    KEY `student_history_info_school_year_semester_teacher_id_index` (`school_year`,`semester`,`teacher_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='学生历史信息';


INSERT INTO `student_history_info` (`school_year`, `semester`, `student_id`, `student_name`, `unit_id`, `unit_name`, `major_code`, `major`, `class_no`, `grade`, `teacher_id`) VALUES ('2012-2013', '1', '13041111', '小明', '04', '电子信息学院（微电子学院）', '0419', '电子信息科学与技术', '17041911', '2017', '13123');
INSERT INTO `student_history_info` (`school_year`, `semester`, `student_id`, `student_name`, `unit_id`, `unit_name`, `major_code`, `major`, `class_no`, `grade`, `teacher_id`) VALUES ('2012-2013', '1', '14184112', '小南', '18', '卓越学院', '1831', '会计学(卓越学院)', '17183111', '2017', NULL);
INSERT INTO `student_history_info` (`school_year`, `semester`, `student_id`, `student_name`, `unit_id`, `unit_name`, `major_code`, `major`, `class_no`, `grade`, `teacher_id`) VALUES ('2012-2013', '1', '15051113', '小贝', '27', '网络空间安全学院（浙江保密学院）', '2724', '网络工程', '17272412', '2017', '31213');
INSERT INTO `student_history_info` (`school_year`, `semester`, `student_id`, `student_name`, `unit_id`, `unit_name`, `major_code`, `major`, `class_no`, `grade`, `teacher_id`) VALUES ('2012-2013', '1', '15141114', '小东', '14', '会计学院', '1406', '会计学', '17140612', '2017', NULL);
INSERT INTO `student_history_info` (`school_year`, `semester`, `student_id`, `student_name`, `unit_id`, `unit_name`, `major_code`, `major`, `class_no`, `grade`, `teacher_id`) VALUES ('2012-2013', '1', '15198115', '小西', '06', '自动化学院（人工智能学院）', '0687', '医学信息工程', '17198711', '2017', '12312');
