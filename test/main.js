import { readFile } from 'fs/promises';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Go の wasm_exec.js を読み込み
const wasmExecPath = join(__dirname, 'wasm_exec.js');
const wasmExecCode = await readFile(wasmExecPath, 'utf-8');

// グローバルスコープに Go を設定
globalThis.Go = eval(wasmExecCode + '; Go');

const runTests = async () => {
    try {
        // WASM ファイルを読み込み
        const wasmPath = join(__dirname, '../main.wasm');
        const wasmBuffer = await readFile(wasmPath);

        // WASM インスタンスを作成
        const go = new Go();
        const wasmModule = await WebAssembly.instantiate(wasmBuffer, go.importObject);

        // WASM を実行
        go.run(wasmModule.instance);

        // テストケース
        const testCases = [
            { input: "48C7C10A00000048C7C3030000004801D94889C8", expected: "d", desc: "Complex calculation" }
        ];

        let allPassed = true;

        for (const [index, test] of testCases.entries()) {
            try {
                console.log(`Running test ${index + 1} (${test.desc})...`);
                const result = globalThis.RunCode(test.input);
                console.log(`Result: ${result}`);
                console.log("-------");

                // エラーチェック
                if (typeof result === 'object' && result.error) {
                    console.error(`❌ Test ${index + 1} (${test.desc}): ERROR - ${result.error}`);
                    allPassed = false;
                    continue;
                }

                const passed = result === test.expected;

                if (passed) {
                    console.log(`✅ Test ${index + 1} (${test.desc}): PASS - got ${result}`);
                } else {
                    console.error(`❌ Test ${index + 1} (${test.desc}): FAIL - Expected ${test.expected}, got ${result}`);
                    allPassed = false;
                }
            } catch (err) {
                console.error(`❌ Test ${index + 1} (${test.desc}): EXCEPTION - ${err.message}`);
                allPassed = false;
            }
        }

        if (allPassed) {
            console.log('\n✅ All tests passed!');
            process.exit(0);
        } else {
            console.error('\n❌ Some tests failed');
            process.exit(1);
        }
    } catch (error) {
        console.error('Test execution error:', error);
        process.exit(1);
    }
};

runTests();
