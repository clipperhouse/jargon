using BenchmarkDotNet.Attributes;
using System;
using System.IO;

namespace Jargon.Benchmark.Benchmarks
{
    public class Tokenize
    {
        private static string Wikipedia;
        [GlobalSetup]
        public void LoadData()
        {
            var path = Path.Combine(Environment.CurrentDirectory, "testdata", "wikipedia.txt");
            Wikipedia = File.ReadAllText(path);
        }

        [Benchmark]
        public void TokenizeBenchmark()
        {
            using (var r = new StringReader(Wikipedia))
            using (var t = new TextTokens(r))
            {
                while (t.MoveNext())
                {
                    // just go, don't use the results
                }
            }
        }
    }
}
