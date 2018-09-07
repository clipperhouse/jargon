using BenchmarkDotNet.Attributes;
using System;
using System.IO;

namespace Jargon.Benchmark.Benchmarks
{
    public class Lemmatizer
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
            var lem = new Jargon.Lemmatizer(dict, 3);
            
            using (var r = new StringReader(Wikipedia))
            using (var tokens = new TextTokens(r))
            using (var l = new LemmaTokens(lem, tokens))
            {
                while (l.MoveNext())
                {
                    // just go, don't use the results
                }
            }
        }
    }
}
