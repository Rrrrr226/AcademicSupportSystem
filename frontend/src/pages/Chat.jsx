import React, { useState, useEffect, useRef } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Layout, Button, List, Typography, message, Spin, Empty, theme, Avatar } from 'antd';
import { PlusOutlined, DeleteOutlined, LeftOutlined, UserOutlined, RobotOutlined } from '@ant-design/icons';
// Assuming @ant-design/x is installed. If not, this line needs the package.
import { Bubble, Sender } from '@ant-design/x'; 
import { chatCompletion, getHistories, getPaginationRecords, delHistory } from '../api/fastgpt';

const { Sider, Content, Header } = Layout;
const { Text, Title } = Typography;

const Chat = () => {
  const { token: { colorBgContainer, colorBorderSecondary } } = theme.useToken();
  const location = useLocation();
  const navigate = useNavigate();
  const { appId, title } = location.state || {}; // Get passed state

  const [histories, setHistories] = useState([]);
  const [activeChatId, setActiveChatId] = useState(null);
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false); // Loading history list
  const [sending, setSending] = useState(false); // Sending message
  const [inputValue, setInputValue] = useState('');

  // Initial Load
  useEffect(() => {
    if (!appId) {
      message.error('缺少应用ID');
      navigate('/subjects');
      return;
    }
    loadHistories();
  }, [appId, navigate]);

  // Load messages when activeChatId changes
  useEffect(() => {
    if (activeChatId) {
      loadMessages(activeChatId);
    } else {
      setMessages([]);
    }
  }, [activeChatId]);

  const loadHistories = async () => {
    setLoading(true);
    try {
      const res = await getHistories({ appId });
      // Adjust according to actual response structure
      // fastgpt usually returns data as array or { list: [] }
      const list = res.data?.data || [];
      setHistories(list);
      
      if (list.length > 0 && !activeChatId) {
        // Automatically select first history? Or start new?
        // Let's create new if empty, or just stay empty
        // setActiveChatId(list[0].chatId);
      }
    } catch (err) {
      console.error(err);
      message.error('加载历史会话失败');
    } finally {
      setLoading(false);
    }
  };

  const loadMessages = async (chatId) => {
    try {
      // Assuming offset=0, pageSize=100 for simplicity
      const res = await getPaginationRecords({ appId, chatId, offset: 0, pageSize: 50 });
      const records = res.data?.data || [];
      // Records might be in reverse order or need formatted
      // FastGPT records usually: { role: 'user'/'assistant', content: '...' }
      // We need to map to Bubble format
      // Map and reverse if needed
      const mapped = records.map(r => ({
        key: r._id || Math.random().toString(),
        role: r.obj === 'Human' || r.role === 'user' ? 'user' : 'ai',
        content: r.value || r.content,
      }));
      // FastGPT histories often return latest first? Check API.
      // Usually chat UI needs oldest first.
      setMessages(mapped.reverse());
    } catch (err) {
      console.error(err);
      message.error('加载消息记录失败');
    }
  };

  const handleNewChat = () => {
    setActiveChatId(null);
    setMessages([]);
  };

  const handleDeleteHistory = async (e, chatId) => {
    e.stopPropagation();
    try {
      await delHistory(appId, chatId);
      message.success('删除成功');
      setHistories(prev => prev.filter(h => h.chatId !== chatId));
      if (activeChatId === chatId) {
        setActiveChatId(null);
      }
    } catch (err) {
      message.error('删除失败');
    }
  };

  const onSend = async (val) => {
    if (!val.trim()) return;
    const currentInput = val;
    setInputValue('');
    setSending(true);

    const newMsg = { key: Date.now().toString(), role: 'user', content: currentInput };
    setMessages(prev => [...prev, newMsg]);

    // Use existing chatId or generate/let backend generate
    // For fastgpt, if we don't pass chatId, it might create one but we need to catch it.
    // Ideally we generate a chatId on frontend for new chat if backend supports it, or use response.
    // Let's rely on backend returning `chatId` or just using the one we have.
    // If activeChatId is null, we generate one or wait for first response?
    // FastGPT API usually accepts `chatId`. 
    const targetChatId = activeChatId || Date.now().toString(); // Simple ID generation
    if (!activeChatId) {
        setActiveChatId(targetChatId);
        // Optimistically add to history list?
        // Better reload histories after first message
    }

    try {
      // Stream handling
      const response = await chatCompletion({
        appId,
        chatId: targetChatId,
        stream: true,
        detail: false,
        messages: [
          ...messages.map(m => ({ 
              role: m.role === 'user' ? 'user' : 'assistant', 
              content: m.content 
          })), // Send full context? Or just last? FastGPT usually handles context if chatId is provided.
          // Wait, FastGPT backend manages history. We typically ONLY send the NEW message if we provide chatId?
          // Check `handler.HandleChatCompletion`. It forwards to FastGPT.
          // FastGPT usually stores history. So we might need to send only the new message or last N messages.
          // Sending ONLY the new message is safer if backend has state.
          { role: 'user', content: currentInput } 
        ]
      });

      // Prepare AI message placeholder
      const aiMsgKey = (Date.now() + 1).toString();
      setMessages(prev => [...prev, { key: aiMsgKey, role: 'ai', content: '' }]);

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let aiContent = '';

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        const chunk = decoder.decode(value);
        // Parse SSE format: data: {...}
        const lines = chunk.split('\n');
        for (const line of lines) {
            if (line.startsWith('data: ')) {
                const jsonStr = line.slice(6);
                if (jsonStr === '[DONE]') continue;
                try {
                    const data = JSON.parse(jsonStr);
                    // FastGPT stream chunk structure: choices[0].delta.content
                    const content = data.choices?.[0]?.delta?.content || '';
                    if (content) {
                        aiContent += content;
                        setMessages(prev => prev.map(m => 
                            m.key === aiMsgKey ? { ...m, content: aiContent } : m
                        ));
                    }
                } catch (e) {
                    // ignore parse error for partial chunks
                }
            }
        }
      }
      
      // Refresh history list if it was a new chat
      if (!histories.find(h => h.chatId === targetChatId)) {
        loadHistories();
      }

    } catch (err) {
      console.error(err);
      message.error('发送失败');
      setMessages(prev => prev.map(m => 
          m.role === 'ai' && m.content === '' ? { ...m, content: 'Error: Failed to get response' } : m
      ));
    } finally {
      setSending(false);
    }
  };

  return (
    <Layout style={{ height: '100vh', background: '#fff' }}>
      <Sider 
        width={300} 
        theme="light" 
        style={{ borderRight: `1px solid ${colorBorderSecondary}` }}
      >
        <div style={{ padding: 16, borderBottom: `1px solid ${colorBorderSecondary}`, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <Button type="text" icon={<LeftOutlined />} onClick={() => navigate('/subjects')}>返回</Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleNewChat}>新对话</Button>
        </div>
        <List
            dataSource={histories}
            loading={loading}
            style={{ height: 'calc(100vh - 64px)', overflowY: 'auto' }}
            renderItem={item => (
                <List.Item 
                    onClick={() => setActiveChatId(item.chatId)}
                    className={activeChatId === item.chatId ? 'ant-list-item-active' : ''}
                    style={{ 
                        cursor: 'pointer', 
                        padding: '12px 16px',
                        background: activeChatId === item.chatId ? '#f0faff' : 'transparent',
                        borderLeft: activeChatId === item.chatId ? '3px solid #1890ff' : '3px solid transparent'
                    }}
                    actions={[
                        <DeleteOutlined onClick={(e) => handleDeleteHistory(e, item.chatId)} style={{ color: '#999' }} />
                    ]}
                >
                    <div style={{ width: '100%', overflow: 'hidden' }}>
                        <Text strong ellipsis>{item.title || '未命名会话'}</Text>
                        <br/>
                        <Text type="secondary" style={{ fontSize: 12 }}>{new Date(item.updateTime || Date.now()).toLocaleDateString()}</Text>
                    </div>
                </List.Item>
            )}
        />
      </Sider>
      
      <Layout>
        <Header style={{ background: colorBgContainer, borderBottom: `1px solid ${colorBorderSecondary}`, padding: '0 24px' }}>
             <Title level={4} style={{ margin: '14px 0' }}>{title || 'AI 助手'}</Title>
        </Header>
        <Content style={{ padding: 24, display: 'flex', flexDirection: 'column' }}>
            <div style={{ flex: 1, overflowY: 'auto', marginBottom: 24 }}>
                {messages.length === 0 ? (
                    <Empty description="开始一个新的对话" />
                ) : (
                    messages.map(msg => (
                        <Bubble
                            key={msg.key}
                            placement={msg.role === 'user' ? 'end' : 'start'}
                            content={msg.content}
                            avatar={<Avatar icon={msg.role === 'user' ? <UserOutlined /> : <RobotOutlined />} />}
                            loading={msg.role === 'ai' && !msg.content && sending}
                        />
                    ))
                )}
            </div>
            <div style={{ maxWidth: 800, margin: '0 auto', width: '100%' }}>
                 <Sender 
                    value={inputValue}
                    onChange={setInputValue}
                    onSubmit={onSend}
                    loading={sending}
                    placeholder="输入您的问题..."
                 />
            </div>
        </Content>
      </Layout>
    </Layout>
  );
};

export default Chat;
