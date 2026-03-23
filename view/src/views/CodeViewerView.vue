<template>
    <div class="code-view-page">
        <header class="hero">
            <div class="hero-copy">
                <p class="eyebrow">即时预览 · Demo</p>
                <h1>查看代码片段的理想界面</h1>
                <p class="lede">
                    支持切换语言与主题，预留操作入口（复制、下载、分享）。当前仅展示前端样式，无后端数据交互。
                </p>
                <div class="hero-actions">
                    <el-button type="primary" @click="informPreview">
                        <el-icon>
                            <View />
                        </el-icon>
                        <span>刷新示例</span>
                    </el-button>
                    <el-button @click="informPreview" plain>
                        <el-icon>
                            <Share />
                        </el-icon>
                        <span>复制链接</span>
                    </el-button>
                </div>
            </div>
            <div class="hero-meta">
                <div class="meta-card">
                    <p class="label">短链</p>
                    <h3>gp-9284a</h3>
                    <p class="muted">示例链接，无真实资源</p>
                </div>
                <div class="meta-card">
                    <p class="label">创建时间</p>
                    <h3>2025-12-26 12:30</h3>
                    <p class="muted">时区按本地浏览器</p>
                </div>
                <div class="meta-card">
                    <p class="label">可见性</p>
                    <h3>公开</h3>
                    <p class="muted">仅 UI 展示</p>
                </div>
            </div>
        </header>

        <el-card class="viewer-card" shadow="hover">
            <div class="toolbar">
                <div class="left">
                    <el-select v-model="language" size="small" class="compact">
                        <el-option label="Go" value="go" />
                        <el-option label="TypeScript" value="ts" />
                        <el-option label="Python" value="py" />
                        <el-option label="JSON" value="json" />
                    </el-select>
                    <el-radio-group v-model="theme" size="small" class="compact">
                        <el-radio-button label="light">浅色</el-radio-button>
                        <el-radio-button label="dark">深色</el-radio-button>
                    </el-radio-group>
                </div>
                <div class="right">
                    <el-button circle text @click="informPreview">
                        <el-icon>
                            <CopyDocument />
                        </el-icon>
                    </el-button>
                    <el-button circle text @click="informPreview">
                        <el-icon>
                            <Download />
                        </el-icon>
                    </el-button>
                    <el-button circle text @click="informPreview">
                        <el-icon>
                            <Share />
                        </el-icon>
                    </el-button>
                </div>
            </div>

            <div class="code-shell" :data-theme="theme">
                <div class="code-header">
                    <span class="dot red"></span>
                    <span class="dot yellow"></span>
                    <span class="dot green"></span>
                    <span class="filename">example.{{ fileExt }}</span>
                </div>
                <div class="code-body">
                    <div class="gutter">
                        <span v-for="(line, idx) in codeLines" :key="idx">{{ idx + 1 }}</span>
                    </div>
                    <pre><code>{{ codeSnippet }}</code></pre>
                </div>
            </div>

            <div class="footer-actions">
                <el-tag size="small" effect="plain">{{ languageLabel }}</el-tag>
                <div class="spacer" />
                <el-button type="primary" plain size="small" @click="informPreview">
                    <el-icon>
                        <CollectionTag />
                    </el-icon>
                    <span>收藏</span>
                </el-button>
                <el-button size="small" @click="informPreview">
                    <el-icon>
                        <Link />
                    </el-icon>
                    <span>生成新短链</span>
                </el-button>
            </div>
        </el-card>
    </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { CollectionTag, CopyDocument, Download, Link, Share, View } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const language = ref<'go' | 'ts' | 'py' | 'json'>('go')
const theme = ref<'light' | 'dark'>('light')

const snippets: Record<string, string> = {
    go: `package main

import "fmt"

func main() {
    fmt.Println("Hello, GopherPaste!")
}
`,
    ts: `export const greet = (name: string) => {
  return ` + '`' + `Hello, ${name}!` + '`' + `
}

console.log(greet('Gopher'))
`,
    py: `from datetime import datetime

print("hello, gopherpaste")
print(datetime.now())
`,
    json: `{
  "name": "gopherpaste",
  "version": "1.0.0",
  "public": true
}
`,
}

const codeSnippet = computed(() => snippets[language.value] || '')
const codeLines = computed(() => codeSnippet.value.split('\n'))
const languageLabel = computed(() => {
    switch (language.value) {
        case 'go':
            return 'Go'
        case 'ts':
            return 'TypeScript'
        case 'py':
            return 'Python'
        case 'json':
            return 'JSON'
        default:
            return 'Plain Text'
    }
})

const fileExt = computed(() => {
    switch (language.value) {
        case 'go':
            return 'go'
        case 'ts':
            return 'ts'
        case 'py':
            return 'py'
        case 'json':
            return 'json'
        default:
            return 'txt'
    }
})

const informPreview = () => {
    ElMessage.info('仅界面展示，功能待接入后端')
}
</script>

<style scoped lang="scss">
.code-view-page {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.hero {
    background: linear-gradient(135deg, #e8f1ff 0%, #f9fbff 50%, #e7f7ff 100%);
    border-radius: 16px;
    padding: 28px;
    display: grid;
    grid-template-columns: 1.6fr 1fr;
    gap: 20px;
    align-items: center;
}

.hero-copy {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .eyebrow {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        padding: 6px 12px;
        background: rgba(64, 158, 255, 0.12);
        color: #2468c0;
        border-radius: 999px;
        font-weight: 700;
        width: fit-content;
    }

    h1 {
        margin: 0;
        font-size: 2.2rem;
        color: #1f2d3d;
    }

    .lede {
        margin: 0;
        color: #596273;
        max-width: 620px;
        line-height: 1.6;
    }

    .hero-actions {
        display: flex;
        gap: 10px;
        flex-wrap: wrap;
    }
}

.hero-meta {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 10px;
}

.meta-card {
    background: #fff;
    border-radius: 12px;
    padding: 14px;
    box-shadow: 0 10px 26px rgba(36, 50, 66, 0.08);

    .label {
        margin: 0;
        color: #909399;
        font-size: 12px;
        letter-spacing: 0.2px;
        text-transform: uppercase;
    }

    h3 {
        margin: 6px 0 4px;
        font-size: 18px;
        color: #1f2d3d;
    }

    .muted {
        margin: 0;
        color: #9aa4b5;
        font-size: 13px;
    }
}

.viewer-card {
    border-radius: 16px;
    border: 1px solid #e4e7ed;
    padding: 14px;

    :deep(.el-card__body) {
        padding: 0;
    }
}

.toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 14px;
    border-bottom: 1px solid #f0f2f5;
    gap: 12px;

    .left {
        display: flex;
        gap: 10px;
        align-items: center;
    }

    .right {
        display: flex;
        gap: 4px;
        align-items: center;
    }

    .compact {
        min-width: 160px;
    }
}

.code-shell {
    margin: 14px;
    border-radius: 12px;
    overflow: hidden;
    border: 1px solid #e4e7ed;
    background: #0f172a;
    color: #e2e8f0;
    transition: background 0.2s ease, color 0.2s ease;

    &[data-theme='light'] {
        background: #f7f9fc;
        color: #1f2d3d;
    }
}

.code-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 14px;
    background: rgba(255, 255, 255, 0.06);
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);

    .dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;
        display: inline-block;
    }

    .red {
        background: #ff5f56;
    }

    .yellow {
        background: #f6bd3b;
    }

    .green {
        background: #51c353;
    }

    .filename {
        margin-left: auto;
        font-family: 'Fira Code', 'Consolas', monospace;
        opacity: 0.8;
    }
}

.code-body {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 0;
    max-height: 520px;
    overflow: hidden;
}

.gutter {
    background: rgba(255, 255, 255, 0.04);
    color: rgba(226, 232, 240, 0.6);
    padding: 14px 12px;
    text-align: right;
    font-family: 'Fira Code', 'Consolas', monospace;
    font-size: 13px;
    min-width: 48px;
    line-height: 1.6;
    border-right: 1px solid rgba(255, 255, 255, 0.08);

    span {
        display: block;
    }
}

pre {
    margin: 0;
    background: transparent;
    color: inherit;
    font-family: 'Fira Code', 'Consolas', monospace;
    padding: 14px 16px;
    overflow: auto;
    white-space: pre;
    line-height: 1.6;
}

.footer-actions {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 14px 14px;
    border-top: 1px solid #f0f2f5;
}

.spacer {
    flex: 1;
}

@media (max-width: 960px) {
    .hero {
        grid-template-columns: 1fr;
    }

    .toolbar {
        flex-direction: column;
        align-items: flex-start;
    }

    .right {
        align-self: flex-end;
    }
}
</style>
