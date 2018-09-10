using BenchmarkDotNet.Attributes;
using System;
using System.IO;

namespace Jargon.Benchmark.Benchmarks
{
    public class TokenizeBench
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
            foreach(var t in Jargon.Tokenize(Wikipedia))
            {
                // left blank
            }
        }
    }
}
