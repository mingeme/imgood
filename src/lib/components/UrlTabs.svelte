<script lang='ts'>
  import { fade } from 'svelte/transition';

  const { urls } = $props<{ urls: string[] }>();
  type Format = 'markdown' | 'html';
  let activeTab = $state<Format>('markdown');

  const codeBlock = $derived(() => {
    if (activeTab === 'markdown') {
      return urls.map((url: string) => {
        const filename = url.split('/').pop();
        return `![${filename}](${url})`;
      }).join('\n');
    }
    else {
      return urls.join('<br>');
    }
  });

  const tabs: Array<{ id: Format; label: string }> = [
    { id: 'markdown', label: 'Markdown' },
    { id: 'html', label: 'HTML' },
  ];
</script>

{#if urls.length > 0}
  <div class='container'>
    <div class='tabs mt-3'>
      <ul>
        {#each tabs as tab, i}
          <li class:is-active={activeTab === tab.id}>
            <a
              href={`#${tab.id}`}
              tabindex={i}
              onclick={() => activeTab = tab.id}
              onkeydown={e => e.key === 'Enter' && (activeTab = tab.id)}
            >
              {tab.label}
            </a>
          </li>
        {/each}
      </ul>
    </div>

    <div transition:fade>
      {#if activeTab === 'markdown'}
        <pre class='code-container'>
          <button
            class='select-all-button'
            title='全选'
            aria-label='全选'
            onclick={() => {
              const codeElement = document.querySelector('.code-block');
              const range = document.createRange();
              range.selectNodeContents(codeElement!!);
              const selection = window.getSelection();
              selection?.removeAllRanges();
              selection?.addRange(range);
            }}
          >
            <svg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'>
              <rect x='9' y='9' width='13' height='13' rx='2' ry='2'></rect>
              <path d='M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1'></path>
            </svg>
          </button>
          <code class='code-block'>{codeBlock()}</code>
        </pre>
      {:else}
        <div class='content'>

        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  pre {
    display: flex;
  }
  .code-container {
    position: relative;
    margin: 0;
    width: 100%;
    overflow: auto;
  }
  .select-all-button {
    position: absolute;
    top: 8px;
    right: 8px;
    padding: 4px;
    background: transparent;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #ccc;
    transition: all 0.2s ease;
    z-index: 10;
  }
  .select-all-button:hover {
    color: #666;
  }
  .select-all-button:focus {
    outline: 2px solid #0066cc;
    outline-offset: 2px;
  }
  .code-block {
    padding: 8px;
    padding-right: 40px;
    width: 100%;
    white-space: pre;
    overflow-x: auto;
  }
</style>
