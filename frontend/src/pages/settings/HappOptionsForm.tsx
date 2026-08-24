import { Divider, Input, Select, Space, Switch } from 'antd';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { SettingListItem } from '@/components/ui';

type Config = Record<string, string | boolean>;

export default function HappOptionsForm({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  const { t } = useTranslation();
  let config: Config = {};
  try { config = JSON.parse(value || '{}') as Config; } catch { config = {}; }
  const set = (key: string, next: string | boolean) => onChange(JSON.stringify({ ...config, [key]: next }));
  const text = (key: string, placeholder = '') => <Input value={String(config[key] ?? '')} placeholder={placeholder} onChange={(e) => set(key, e.target.value)} />;
  const toggle = (key: string) => <Switch checked={Boolean(config[key])} onChange={(v) => set(key, v)} />;
  const row = (key: string, control: ReactNode) => <SettingListItem paddings="small" title={t(`pages.settings.happ.${key}`)}>{control}</SettingListItem>;

  return <Space direction="vertical" style={{ width: '100%' }} size={0}>
    <Divider orientation="left">{t('pages.settings.happ.announcement')}</Divider>
    {row('infoColor', <Select style={{ width: 180 }} value={String(config.infoColor ?? '') || undefined} allowClear options={['blue','red','green'].map(v => ({ value: v, label: v }))} onChange={(v) => set('infoColor', v ?? '')} />)}
    {row('infoText', text('infoText'))}{row('infoButtonText', text('infoButtonText'))}{row('infoButtonLink', text('infoButtonLink', 'https://...'))}
    <Divider orientation="left">{t('pages.settings.happ.subscription')}</Divider>
    {row('expireEnable', toggle('expireEnable'))}{row('expireButtonLink', text('expireButtonLink', 'https://...'))}
    {row('hideServerSettings', toggle('hideServerSettings'))}{row('autoUpdateOnOpen', toggle('autoUpdateOnOpen'))}
    {row('pinCurrent', toggle('pinCurrent'))}{row('autoUpdate', toggle('autoUpdate'))}
    <Divider orientation="left">{t('pages.settings.happ.fragmentation')}</Divider>
    {row('fragmentationEnable', toggle('fragmentationEnable'))}{row('fragmentationPackets', text('fragmentationPackets', 'tlshello'))}
    {row('fragmentationLength', text('fragmentationLength', '50-100'))}{row('fragmentationInterval', text('fragmentationInterval', '10-20'))}
    {row('fragmentationMaxSplit', text('fragmentationMaxSplit', '100-200'))}
    <Divider orientation="left">{t('pages.settings.happ.noises')}</Divider>
    {row('noisesEnable', toggle('noisesEnable'))}
    {row('noisesType', <Select style={{ width: 180 }} value={String(config.noisesType ?? '') || undefined} allowClear options={['array','str','hex','base64'].map(v => ({ value: v, label: v }))} onChange={(v) => set('noisesType', v ?? '')} />)}
    {row('noisesPacket', text('noisesPacket'))}{row('noisesDelay', text('noisesDelay', '50'))}{row('noisesApplyTo', text('noisesApplyTo'))}
    <Divider orientation="left">Ping</Divider>
    {row('pingType', <Select style={{ width: 180 }} value={String(config.pingType ?? '') || undefined} allowClear options={['proxy','proxy-head','tcp','icmp'].map(v => ({ value: v, label: v }))} onChange={(v) => set('pingType', v ?? '')} />)}
    {row('pingResult', <Select style={{ width: 180 }} value={String(config.pingResult ?? '') || undefined} allowClear options={['time','speed'].map(v => ({ value: v, label: v }))} onChange={(v) => set('pingResult', v ?? '')} />)}
    {row('subscriptionSort', <Select style={{ width: 180 }} value={String(config.subscriptionSort ?? '') || undefined} allowClear options={['ping','name','country'].map(v => ({ value: v, label: v }))} onChange={(v) => set('subscriptionSort', v ?? '')} />)}
  </Space>;
}
