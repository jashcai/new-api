/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useContext, useEffect, useState } from 'react';
import {
  Button,
  Typography,
  Input,
  ScrollList,
  ScrollItem,
} from '@douyinfe/semi-ui';
import { API, showError, copy, showSuccess } from '../../helpers';
import { useIsMobile } from '../../hooks/common/useIsMobile';
import { API_ENDPOINTS } from '../../constants/common.constant';
import { StatusContext } from '../../context/Status';
import { useActualTheme } from '../../context/Theme';
import { marked } from 'marked';
import { useTranslation } from 'react-i18next';
import {
  IconGithubLogo,
  IconPlay,
  IconFile,
  IconCopy,
} from '@douyinfe/semi-icons';
import { Link } from 'react-router-dom';
import NoticeModal from '../../components/layout/NoticeModal';
import {
  Activity,
  Bot,
  BrainCircuit,
  CheckCircle2,
  FileCheck2,
  Gauge,
  KeyRound,
  LockKeyhole,
  Network,
  Route,
  ShieldCheck,
  Sparkles,
} from 'lucide-react';

const { Text } = Typography;

const Home = () => {
  const { t, i18n } = useTranslation();
  const [statusState] = useContext(StatusContext);
  const actualTheme = useActualTheme();
  const [homePageContentLoaded, setHomePageContentLoaded] = useState(false);
  const [homePageContent, setHomePageContent] = useState('');
  const [noticeVisible, setNoticeVisible] = useState(false);
  const isMobile = useIsMobile();
  const isDemoSiteMode = statusState?.status?.demo_site_enabled || false;
  const docsLink = statusState?.status?.docs_link || '';
  const serverAddress =
    statusState?.status?.server_address || `${window.location.origin}`;
  const endpointItems = API_ENDPOINTS.map((e) => ({ value: e }));
  const [endpointIndex, setEndpointIndex] = useState(0);
  const modelIcons = [
    Network,
    Bot,
    BrainCircuit,
    Sparkles,
    Route,
    KeyRound,
    Activity,
    ShieldCheck,
    Gauge,
    LockKeyhole,
    FileCheck2,
    CheckCircle2,
  ];
  const capabilityCards = [
    {
      key: 'models',
      title: '统一模型入口',
      desc: '通过一个标准接口聚合企业内外部模型，应用只需要接入一次，就能按权限访问已批准的模型能力。',
      action: '浏览模型',
    },
    {
      key: 'availability',
      title: '更高可用路由',
      desc: '根据渠道健康度、延迟和策略自动分配请求，在供应商波动时切换备用通道，减少业务中断。',
      action: '查看路由',
    },
    {
      key: 'performance',
      title: '成本与性能平衡',
      desc: '持续观察吞吐、延迟、错误率和用量成本，在体验、速度和预算之间做可控取舍。',
      action: '查看看板',
    },
    {
      key: 'policy',
      title: '数据与审计策略',
      desc: '用团队、模型、额度和审计规则约束访问范围，确保密钥、日志和调用轨迹留在企业控制内。',
      action: '查看策略',
    },
  ];

  const renderCapabilityVisual = (cardKey) => {
    if (cardKey === 'models') {
      return (
        <div className='classic-capability-icon-cloud'>
          {modelIcons.map((Icon, index) => (
            <span className='classic-capability-icon-node' key={index}>
              <Icon size={15} strokeWidth={1.8} />
            </span>
          ))}
        </div>
      );
    }

    if (cardKey === 'availability') {
      return (
        <div className='classic-capability-route-map'>
          <span className='classic-capability-model-chip'>
            anthropic/claude-opus-4.8
          </span>
          <div className='classic-capability-route-lines'>
            <span />
            <span />
            <span />
          </div>
          <div className='classic-capability-route-nodes'>
            <span>
              <Route size={18} />
            </span>
            <span>
              <Bot size={18} />
            </span>
            <span>
              <Network size={18} />
            </span>
          </div>
        </div>
      );
    }

    if (cardKey === 'performance') {
      return (
        <div className='classic-capability-chart-stack'>
          <div className='classic-capability-chart-card back'>
            <span>{t('吞吐')}</span>
            <div className='classic-mini-chart teal' />
          </div>
          <div className='classic-capability-chart-card front'>
            <span>{t('延迟')}</span>
            <div className='classic-mini-chart amber' />
            <div className='classic-mini-chart blue' />
          </div>
        </div>
      );
    }

    return (
      <div className='classic-capability-policy-visual'>
        <div className='classic-policy-locks'>
          <LockKeyhole size={16} />
          <CheckCircle2 size={30} />
          <LockKeyhole size={16} />
        </div>
        <ShieldCheck size={94} strokeWidth={1.1} />
      </div>
    );
  };

  const displayHomePageContent = async () => {
    setHomePageContent(localStorage.getItem('home_page_content') || '');
    const res = await API.get('/api/home_page_content');
    const { success, message, data } = res.data;
    if (success) {
      let content = data;
      if (!data.startsWith('https://')) {
        content = marked.parse(data);
      }
      setHomePageContent(content);
      localStorage.setItem('home_page_content', content);

      // 如果内容是 URL，则发送主题模式
      if (data.startsWith('https://')) {
        const iframe = document.querySelector('iframe');
        if (iframe) {
          iframe.onload = () => {
            iframe.contentWindow.postMessage({ themeMode: actualTheme }, '*');
            iframe.contentWindow.postMessage({ lang: i18n.language }, '*');
          };
        }
      }
    } else {
      showError(message);
      setHomePageContent('加载首页内容失败...');
    }
    setHomePageContentLoaded(true);
  };

  const handleCopyBaseURL = async () => {
    const ok = await copy(serverAddress);
    if (ok) {
      showSuccess(t('已复制到剪切板'));
    }
  };

  useEffect(() => {
    const checkNoticeAndShow = async () => {
      const lastCloseDate = localStorage.getItem('notice_close_date');
      const today = new Date().toDateString();
      if (lastCloseDate !== today) {
        try {
          const res = await API.get('/api/notice');
          const { success, data } = res.data;
          if (success && data && data.trim() !== '') {
            setNoticeVisible(true);
          }
        } catch (error) {
          console.error('获取公告失败:', error);
        }
      }
    };

    checkNoticeAndShow();
  }, []);

  useEffect(() => {
    displayHomePageContent().then();
  }, []);

  useEffect(() => {
    const timer = setInterval(() => {
      setEndpointIndex((prev) => (prev + 1) % endpointItems.length);
    }, 3000);
    return () => clearInterval(timer);
  }, [endpointItems.length]);

  return (
    <div className='classic-page-fill classic-home-page w-full overflow-x-hidden'>
      <NoticeModal
        visible={noticeVisible}
        onClose={() => setNoticeVisible(false)}
        isMobile={isMobile}
      />
      {homePageContentLoaded && homePageContent === '' ? (
        <div className='classic-home-default w-full overflow-x-hidden'>
          {/* Banner 部分 */}
          <div className='classic-home-hero w-full border-b border-semi-color-border relative overflow-x-hidden'>
            {/* 背景模糊晕染球 */}
            <div className='blur-ball blur-ball-indigo' />
            <div className='blur-ball blur-ball-teal' />
            <div className='classic-home-shell px-4 pt-24 pb-12 md:pt-28 md:pb-16'>
              <div className='classic-home-grid mx-auto max-w-6xl'>
                <div className='classic-home-copy'>
                  <h1
                    className='text-4xl md:text-5xl lg:text-6xl xl:text-7xl font-bold text-semi-color-text-0 leading-tight'
                  >
                    <>
                      <span className='shine-text'>
                        ONE API,
                        <br />
                        ONE PLATFORM
                      </span>
                    </>
                  </h1>
                  <p className='text-base md:text-lg text-semi-color-text-1 mt-4 md:mt-6 max-w-2xl leading-relaxed'>
                    {t(
                      '通过一个受控入口路由 OpenAI 兼容、Claude、Gemini 和私有模型流量，集中管理密钥、额度、计费、审计日志和团队策略。',
                    )}
                  </p>
                  {/* BASE URL 与端点选择 */}
                  <div className='flex flex-col md:flex-row items-center justify-start gap-4 w-full mt-6 md:mt-8 max-w-xl'>
                    <Input
                      readonly
                      value={serverAddress}
                      className='classic-endpoint-input flex-1 !rounded-full'
                      size={isMobile ? 'default' : 'large'}
                      suffix={
                        <div className='flex items-center gap-2'>
                          <ScrollList
                            bodyHeight={32}
                            style={{ border: 'unset', boxShadow: 'unset' }}
                          >
                            <ScrollItem
                              mode='wheel'
                              cycled={true}
                              list={endpointItems}
                              selectedIndex={endpointIndex}
                              onSelect={({ index }) => setEndpointIndex(index)}
                            />
                          </ScrollList>
                          <Button
                            type='primary'
                            onClick={handleCopyBaseURL}
                            icon={<IconCopy />}
                            className='!rounded-full'
                          />
                        </div>
                      }
                    />
                  </div>

                {/* 操作按钮 */}
                <div className='flex flex-row gap-4 justify-start items-center mt-6'>
                  <Link to='/console'>
                    <Button
                      theme='solid'
                      type='primary'
                      size={isMobile ? 'default' : 'large'}
                      className='classic-primary-cta !rounded-3xl px-8 py-2'
                      icon={<IconPlay />}
                    >
                      {t('进入控制台')}
                    </Button>
                  </Link>
                  {isDemoSiteMode && statusState?.status?.version ? (
                    <Button
                      size={isMobile ? 'default' : 'large'}
                      className='flex items-center !rounded-3xl px-6 py-2'
                      icon={<IconGithubLogo />}
                      onClick={() =>
                        window.open(
                          'https://github.com/QuantumNous/new-api',
                          '_blank',
                        )
                      }
                    >
                      {statusState.status.version}
                    </Button>
                  ) : (
                    docsLink && (
                      <Button
                        size={isMobile ? 'default' : 'large'}
                        className='flex items-center !rounded-3xl px-6 py-2'
                        icon={<IconFile />}
                        onClick={() => window.open(docsLink, '_blank')}
                      >
                        {t('文档')}
                      </Button>
                    )
                  )}
                </div>
                </div>

                <div className='classic-gateway-panel'>
                  <div className='classic-gateway-panel-header classic-gateway-visual-banner'>
                    <span className='classic-gateway-banner-status'>
                      {t('策略已启用')}
                    </span>
                  </div>
                  <div className='classic-gateway-route'>
                    <div>
                      <div className='classic-gateway-label'>
                        {t('当前入口')}
                      </div>
                      <div className='classic-gateway-value'>
                        {serverAddress}
                      </div>
                    </div>
                    <div className='classic-gateway-status'>200 OK</div>
                  </div>
                  <div className='classic-gateway-metrics'>
                    {[
                      ['40+', '上游供应商'],
                      ['100+', '策略治理模型'],
                      ['24/7', '运营可见性'],
                    ].map(([value, label]) => (
                      <div className='classic-gateway-metric' key={label}>
                        <div className='classic-gateway-metric-value'>
                          {value}
                        </div>
                        <div className='classic-gateway-metric-label'>
                          {t(label)}
                        </div>
                      </div>
                    ))}
                  </div>
                  <div className='classic-gateway-flow'>
                    {[
                      ['应用请求', '/v1/chat/completions'],
                      ['策略校验', 'group: default'],
                      ['渠道路由', 'latency-aware'],
                      ['审计记录', 'usage + cost'],
                    ].map(([title, desc]) => (
                      <div className='classic-gateway-flow-row' key={title}>
                        <span className='classic-gateway-flow-dot' />
                        <div>
                          <div className='font-semibold text-semi-color-text-0'>
                            {t(title)}
                          </div>
                          <div className='text-xs text-semi-color-text-2 mt-0.5'>
                            {desc}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                <div className='classic-home-capabilities'>
                  <div className='flex items-center mb-6 justify-between'>
                    <Text
                      type='tertiary'
                      className='text-base md:text-lg font-medium tracking-wide'
                    >
                      {t('企业运营层')}
                    </Text>
                    <Text type='tertiary' className='text-sm'>
                      {t('统一接入 · 集中治理 · 可审计运营')}
                    </Text>
                  </div>
                  <div className='classic-enterprise-grid'>
                    {capabilityCards.map((card) => (
                      <div className='classic-enterprise-card' key={card.key}>
                        <div className='classic-capability-visual'>
                          {renderCapabilityVisual(card.key)}
                        </div>
                        <div className='classic-capability-copy'>
                          <div className='classic-capability-title'>
                            {t(card.title)}
                          </div>
                          <div className='classic-capability-desc'>
                            {t(card.desc)}
                          </div>
                          <Link to='/pricing' className='classic-capability-link'>
                            {t(card.action)}
                          </Link>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      ) : (
        <div className='classic-page-fill overflow-x-hidden w-full'>
          {homePageContent.startsWith('https://') ? (
            <iframe
              src={homePageContent}
              className='w-full h-screen border-none'
            />
          ) : (
            <div
              className='mt-[60px]'
              dangerouslySetInnerHTML={{ __html: homePageContent }}
            />
          )}
        </div>
      )}
    </div>
  );
};

export default Home;
