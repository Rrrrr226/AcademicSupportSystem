import React, { useState } from 'react';
import { Card, Button, Upload, message, Typography, Statistic, Row, Col, Alert } from 'antd';
import { UploadOutlined, FileExcelOutlined, DownloadOutlined } from '@ant-design/icons';
import { importStudentSubjects, downloadImportTemplate } from '../../api';

const { Title } = Typography;
const { Dragger } = Upload;

const ImportTab = () => {
  const [importResult, setImportResult] = useState(null);

  const handleImportExcel = async (file) => {
    const token = localStorage.getItem('adminToken');
    setImportResult(null);
    try {
      const response = await importStudentSubjects(file, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        setImportResult(response.data.data);
        message.success('导入完成');
      } else {
        message.error(response.data?.message || '导入失败');
      }
    } catch (error) {
      console.error('导入失败:', error);
      message.error(error.response?.data?.message || '导入失败');
    }
    return false; // 阻止默认上传行为
  };

  const handleDownloadTemplate = async () => {
    try {
      const token = localStorage.getItem('adminToken');
      const response = await downloadImportTemplate(token);
      
      // 创建 Blob URL 并下载
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'subject_import_template.xlsx');
      document.body.appendChild(link);
      link.click();
      
      // 清理
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('下载模板失败:', error);
      message.error('下载模板失败');
    }
  };

  const uploadProps = {
    name: 'file',
    multiple: false,
    accept: '.xlsx,.xls',
    beforeUpload: handleImportExcel,
    showUploadList: false,
  };

  return (
    <div>
      <Title level={4}>导入学生科目名单</Title>
      <Alert
        message="Excel 格式要求"
        description={
          <div>
            <p>Excel 文件需要包含以下列：</p>
            <ul>
              <li><strong>学号</strong>（或 staff_id、StaffId）- 学生学号</li>
              <li><strong>科目名称</strong>（或 科目、subject_name、SubjectName）- 科目名称（需要与系统中已有科目名称一致）</li>
            </ul>
            <p>每行代表一个学生-科目的对应关系，同一学生可以有多行对应不同科目。</p>
          </div>
        }
        type="info"
        showIcon
        style={{ marginBottom: 16 }}
      />

      <div style={{ marginBottom: 24 }}>
        <Button 
          type="default" 
          icon={<DownloadOutlined />}
          onClick={handleDownloadTemplate}
        >
          下载导入模板
        </Button>
      </div>

      <Dragger {...uploadProps} style={{ marginBottom: 24 }}>
        <p className="ant-upload-drag-icon">
          <UploadOutlined style={{ fontSize: 48, color: '#1890ff' }} />
        </p>
        <p className="ant-upload-text">点击或拖拽 Excel 文件到此区域上传</p>
        <p className="ant-upload-hint">
          支持 .xlsx 和 .xls 格式
        </p>
      </Dragger>

      {importResult && (
        <Card title="导入结果" style={{ marginTop: 24 }}>
          <Row gutter={16}>
            <Col span={8}>
              <Statistic title="总记录数" value={importResult.total} />
            </Col>
            <Col span={8}>
              <Statistic
                title="成功数"
                value={importResult.successCount}
                valueStyle={{ color: '#3f8600' }}
              />
            </Col>
            <Col span={8}>
              <Statistic
                title="失败数"
                value={importResult.failCount}
                valueStyle={{ color: importResult.failCount > 0 ? '#cf1322' : undefined }}
              />
            </Col>
          </Row>
          {importResult.errors && importResult.errors.length > 0 && (
            <div style={{ marginTop: 16 }}>
              <Title level={5}>错误详情：</Title>
              <ul>
                {importResult.errors.slice(0, 10).map((err, idx) => (
                  <li key={idx} style={{ color: '#cf1322' }}>{err}</li>
                ))}
                {importResult.errors.length > 10 && (
                  <li>...还有 {importResult.errors.length - 10} 个错误</li>
                )}
              </ul>
            </div>
          )}
        </Card>
      )}
    </div>
  );
};

export default ImportTab;
