using CommandLine;
using System;
using System.Diagnostics;
using System.IO;

namespace Jargon.Cmd
{
    public sealed class Program
    {
        private sealed class Options
        {
            [Option('f', Required = false, HelpText = "A file path to lemmatize")]
            public string File { get; set; }

            internal bool FileIsSet => !string.IsNullOrEmpty(File);

            [Option('s', Required = false, HelpText = "A (quoted) string to lemmatize")]
            public string String { get; set; }

            internal bool StringIsSet => !string.IsNullOrEmpty(String);

            [Option('u', Required = false, HelpText = "A URL to fetch and lemmatize")]
            public string Url { get; set; }

            internal bool UrlIsSet => !string.IsNullOrEmpty(Url);
        }

        public static void Main(string[] args)
        {
            var res = Parser.Default.ParseArguments<Options>(args);
            res.WithParsed(HandleOptions).WithNotParsed(e => PrintExamplesAndExit());
        }

        private static void HandleOptions(Options opts)
        {
            var numSet =
                (opts.FileIsSet ? 1 : 0) +
                (opts.StringIsSet ? 1 : 0) +
                (opts.UrlIsSet ? 1 : 0);
            if (numSet > 1)
            {
                Console.WriteLine($"Only one of `f`, `s`, and `u` may be set");
                Environment.Exit(-2);
            }

            if (opts.FileIsSet)
            {
                if (!File.Exists(opts.File))
                {
                    Console.WriteLine($"Could not find file: {opts.File}");
                    Environment.Exit(-3);
                }

                try
                {
                    using (var reader = new StreamReader(File.OpenRead(opts.File)))
                    {
                        Lemmatize(reader);
                        return;
                    }
                }
                catch(Exception e)
                {
                    Console.WriteLine($"Could not read file ({e.Message}): {opts.File}");
                    Environment.Exit(-4);
                }   
            }

            if (opts.StringIsSet)
            {
                using (var reader = new StringReader(opts.String))
                {
                    Lemmatize(reader);
                    return;
                }
            }

            if (opts.UrlIsSet)
            {
                try
                {
                    using (var web = new CompressedWebClient())
                    {
                        var html = web.DownloadString(opts.Url);
                        using (var reader = new StringReader(html))
                        {
                            Lemmatize(reader);
                        }
                        return;
                    }
                }
                catch (Exception e)
                {
                    Console.WriteLine($"Could not download url ({e.Message}): {opts.Url}");
                    Environment.Exit(-5);
                }
            }

            Lemmatize(Console.In);
        }

        private static void Lemmatize(TextReader reader)
        {
            var lemmatizer = new Lemmatizer(Data.StackExchange.Instance, 3);
            using(var toks = new TextTokens(reader))
            using(var e = new LemmaTokens(in lemmatizer, toks))
            {
                while (e.MoveNext())
                {
                    Console.Write(e.Current.Value);
                }
            }
        }

        private static void PrintExamplesAndExit()
        {
            string cmdName;
            using (var proc = Process.GetCurrentProcess())
            {
                cmdName = proc.ProcessName;
            }

            var pathSeparator = Path.DirectorySeparatorChar;

            Console.WriteLine($@"
Usage: {cmdName} accepts piped UTF8 text from tools such as cat, curl or echo, via Stdin
		
  Example: echo ""I luv Rails"" | {cmdName}
Alternatively, use {cmdName} 'standalone' by passing flags for text sources:
  Example: %s -f {pathSeparator}path{pathSeparator}to{pathSeparator}file.txt

Results are piped to Stdout (regardless of input)");

            Environment.Exit(-1);
        }
    }
}
