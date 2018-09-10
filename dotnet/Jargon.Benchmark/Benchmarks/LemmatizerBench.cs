using BenchmarkDotNet.Attributes;
using System;
using System.IO;

namespace Jargon.Benchmark.Benchmarks
{
    public class LemmatizerBench
    {
        private static string Wikipedia;
        [GlobalSetup]
        public void LoadData()
        {
            var path = Path.Combine(Environment.CurrentDirectory, "testdata", "wikipedia.txt");
            Wikipedia = File.ReadAllText(path);
        }

        [Benchmark]
        public void LemmatizerBenchmark()
        {
            var dict = Data.StackExchange.Instance;
            var lem = new Lemmatizer(dict, 3);

            foreach (var _ in lem.Lemmatize(Wikipedia))
            {
                // left blank
            }
        }
    }
}
