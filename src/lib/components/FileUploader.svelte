<script lang='ts'>
  import type { FileWithPreview } from '$lib/types';
  import { deserialize } from '$app/forms';
  import { calculateFileHash } from '$lib/utils/file';
  import { onDestroy, onMount } from 'svelte';

  const { user, files, onFilesSelected, onUploadSuccess, onRemove, onClear } = $props<{
    user: any;
    files: FileWithPreview[];
    onFilesSelected: (files: FileWithPreview[]) => void;
    onUploadSuccess: (url: string) => void;
    onRemove: (index: number) => void;
    onClear: () => void;
  }>();

  let dragCounter = $state(0);

  function handleGlobalPaste(e: ClipboardEvent) {
    if (!user)
      return;
    const pastedFiles = Array.from(e.clipboardData?.files ?? [])
      .filter(file => file.type.startsWith('image/'));
    if (pastedFiles.length > 0) {
      getUniqueFiles(pastedFiles).then((uniqueFiles) => {
        if (uniqueFiles.length > 0) {
          onFilesSelected(uniqueFiles);
        }
      });
    }
  }

  onMount(() => {
    if (import.meta.env.SSR)
      return;
    document.addEventListener('paste', handleGlobalPaste);
  });

  onDestroy(() => {
    if (import.meta.env.SSR)
      return;
    document.removeEventListener('paste', handleGlobalPaste);
  });

  function handleDragEnter(e: DragEvent) {
    e.preventDefault();
    dragCounter++;
  }

  function handleDragLeave(e: DragEvent) {
    e.preventDefault();
    dragCounter--;
  }

  async function createFileWithPreview(file: File): Promise<FileWithPreview> {
    if (!file || !(file instanceof File)) {
      throw new Error('Invalid file');
    }
    const hash = await calculateFileHash(file);
    const preview = URL.createObjectURL(file);

    return { file, preview, hash, status: 'pending' };
  }

  async function getUniqueFiles(filesArray: File[]): Promise<FileWithPreview[]> {
    const seenHashes = new Set<string>(files.map((existingFile: FileWithPreview) => existingFile.hash));

    const filesWithHash = await Promise.all(filesArray.map(async (file) => {
      if (file instanceof File) {
        return createFileWithPreview(file);
      }
      else {
        console.error('Invalid file:', file);
        return null;
      }
    }));

    return filesWithHash.filter((newFile): newFile is FileWithPreview => {
      if (newFile === null)
        return false;

      if (seenHashes.has(newFile.hash)) {
        return false;
      }

      seenHashes.add(newFile.hash);
      return true;
    });
  }

  async function handleDrop(e: DragEvent) {
    e.preventDefault();
    dragCounter = 0;
    if (!user)
      return;

    const droppedFiles = Array.from(e.dataTransfer?.files ?? [])
      .filter(file => file.type.startsWith('image/'));
    if (droppedFiles.length > 0) {
      const uniqueFiles = await getUniqueFiles(droppedFiles);
      if (uniqueFiles.length > 0) {
        onFilesSelected(uniqueFiles);
      }
    }
  }

  async function upload() {
    for (const fileWithPreview of files) {
      if (fileWithPreview.status !== 'pending')
        continue;
      fileWithPreview.status = 'uploading';
      const hash = await calculateFileHash(fileWithPreview.file);
      const res = await fetch('?/upload', {
        method: 'POST',
        body: JSON.stringify({
          name: fileWithPreview.file.name,
          size: fileWithPreview.file.size,
          hash,
        }),
      });
      const result = deserialize(await res.text());
      const { type } = result;

      if (type === 'success') {
        const { signedUrl, url } = result.data as { signedUrl: string; url: string };
        await fetch(signedUrl, {
          method: 'PUT',
          body: fileWithPreview.file,
        });
        fileWithPreview.status = 'success';
        onUploadSuccess(url);
      }
      else if (type === 'failure') {
        fileWithPreview.status = 'duplicated';
      }
      else if (type === 'error') {
        fileWithPreview.status = 'error';
      }
    }
  }
</script>

<div>
  <div class='file is-boxed file-cta image-upload-height'>
    <div
      class='container file-label is-responsive'
      class:cursor-default={!user}
      class:is-dragover={dragCounter > 0}
      role='button'
      tabindex='0'
      ondragenter={handleDragEnter}
      ondragleave={handleDragLeave}
      ondragover={e => e.preventDefault()}
      ondrop={handleDrop}
    >
      <span class='image-upload-content'>
        {#if !user}
          <div class='has-text-centered'>
            You need to <a href='/signin'>sign in</a> before you can upload images
          </div>
        {:else if files.length > 0}
          {#each files as { preview, file, status }, i}
            <div class='preview-container'>
              <button
                type='button'
                class='delete-button'
                onclick={(e) => {
                  e.stopPropagation();
                  onRemove(i);
                }}
                aria-label='Delete image'
              >
                ×
              </button>
              <img
                src={preview}
                alt={file.name}
                class='preview-image'
              />
              <div class='file-info'>
                <span class='file-name is-small' title={file.name}>
                  {file.name}
                </span>
                {#if status !== 'pending'}
                  <div class='progress-container'>
                    <progress
                      class='progress is-small'
                      class:is-success={status === 'success'}
                      class:is-danger={status === 'error' || status === 'duplicated'}
                      class:is-primary={status === 'uploading'}
                      max='100'
                      value={status === 'uploading' ? 20 : 100}
                    ></progress>
                    <span class='status-text'
                          class:has-text-success={status === 'success'}
                          class:has-text-danger={status === 'error' || status === 'duplicated'}>
                      {#if status === 'uploading'}
                        Uploading...
                      {:else if status === 'success'}
                        Upload completed
                      {:else if status === 'error'}
                        Upload failed
                      {:else if status === 'duplicated'}
                        Already exists
                      {/if}
                    </span>
                  </div>
                {/if}
              </div>
            </div>
          {/each}
        {:else}
          <div class='has-text-centered'>
            Drag &amp; drop files here ...
            <br />
            or
            <br />
            Copy &amp; paste screenshots here ...
          </div>
        {/if}
      </span>
    </div>
  </div>

  {#if user}
    <div class='is-flex'>
      <label class='file-label container select-files-label'>
        <input
          id='file-upload'
          class='file-input select-files-input'
          type='file'
          name='image'
          accept='.jpg, .jpeg, .png, .gif, .bmp'
          multiple
          onchange={async (e: Event) => {
            const target = e.target as HTMLInputElement;
            const filesArray = Array.from(target.files ?? []);
            const uniqueFiles = await getUniqueFiles(filesArray);
            if (uniqueFiles.length > 0) {
              onFilesSelected(uniqueFiles);
            }
          }}
        />
        <span class='select-files'>Select images...</span>
      </label>

      {#if files.length > 0}
        <div class='is-flex'>
          <button
            class='button is-danger select-button'
            type='button'
            onclick={onClear}
          >
            Clear
          </button>
          <button
            class='button is-primary select-button'
            onclick={upload}
          >
            Upload
          </button>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .preview-container {
    max-width: 200px;
    position: relative;
    display: inline-block;
    margin: 0.5rem;
    border-radius: 0.25rem;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .delete-button {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    width: 1.5rem;
    height: 1.5rem;
    border-radius: 50%;
    background: rgba(0, 0, 0, 0.5);
    color: white;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1rem;
    transition: background-color 0.2s;
  }

  .delete-button:hover {
    background: rgba(0, 0, 0, 0.7);
  }

  .preview-image {
    width: 200px;
    height: 200px;
    object-fit: cover;
    display: block;
  }

  .image-upload-height {
    height: 100%;
    min-height: 20rem;
    border: 2px dashed var(--bulma-border);
    border-radius: 0.5rem;
    transition: all 0.2s ease;
  }

  .image-upload-height:hover,
  .is-dragover {
    border-color: var(--bulma-primary);
    background-color: rgba(var(--bulma-primary-rgb), 0.05);
  }

  .image-upload-content {
    margin-block: auto;
    overflow: auto;
    padding: 1rem;
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    justify-content: center;
  }

  .select-files-label {
    border: 1px solid var(--bulma-border);
    border-radius: 0.25rem;
    overflow: hidden;
  }

  .select-files {
    display: flex;
    align-items: center;
    padding: 0.5rem 1rem;
  }

  .select-files-input {
    cursor: pointer;
  }

  .cursor-default {
    cursor: default;
  }

  .select-button {
    margin-left: 0.5rem;
  }

  .progress-container {
    position: relative;
    margin-top: 0.5rem;
  }

  .status-text {
    left: 0;
    right: 0;
    text-align: center;
    font-size: 0.75rem;
    margin-top: 0.25rem;
    color: #4a4a4a;
  }

  .progress {
    margin-top: 0.5rem;
    margin-bottom: 0.5rem;
    height: 4px;
  }

  .progress::-webkit-progress-value {
    transition: width 0.3s ease;
  }

</style>
