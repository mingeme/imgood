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
        <pre>
        <code>{codeBlock()}</code>
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
</style>
