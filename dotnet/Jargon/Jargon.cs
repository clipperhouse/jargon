using Jargon.Impl;
using System;
using System.Collections.Generic;
using System.IO;

namespace Jargon
{
    public static class Jargon
    {
        public static IEnumerable<Token> Tokenize(string str)
        {
            if (str == null) throw new ArgumentNullException(nameof(str));

            var strReader = new StringReader(str);
            return Tokenize(strReader);
        }

        public static IEnumerable<Token> Tokenize(TextReader reader)
        {
            if(reader == null) throw new ArgumentNullException(nameof(reader));

            return new ITokensToEnumerableAdapter(new TextTokens(reader));
        }

        public static IEnumerable<Token> TokenizeHTML(string str)
        {
            if (str == null) throw new ArgumentNullException(nameof(str));

            var strReader = new StringReader(str);
            return TokenizeHTML(strReader);
        }

        public static IEnumerable<Token> TokenizeHTML(TextReader reader)
        {
            if (reader == null) throw new ArgumentNullException(nameof(reader));

            return new ITokensToEnumerableAdapter(new HTMLTokens(reader));
        }

        public static IEnumerable<Token> Lemmatize(string str, Lemmatizer lem)
        {
            if (str == null) throw new ArgumentNullException(nameof(str));

            var strReader = new StringReader(str);
            return Lemmatize(strReader, lem);
        }

        public static IEnumerable<Token> Lemmatize(TextReader reader, Lemmatizer lem)
        {
            if (reader == null) throw new ArgumentNullException(nameof(reader));

            var tokens = Tokenize(reader);
            return Lemmatize(tokens, lem);
        }

        public static IEnumerable<Token> Lemmatize(IEnumerable<Token> tokens, Lemmatizer lem)
        {
            if (tokens == null) throw new ArgumentNullException(nameof(tokens));

            return 
                new ITokensToEnumerableAdapter(
                    new LemmaTokens<EnumerableToITokensAdapter>(
                        lem, 
                        new EnumerableToITokensAdapter(tokens.GetEnumerator())
                    )
                );
        }

        public static IEnumerable<Token> LemmatizeHTML(string str, Lemmatizer lem)
        {
            if (str == null) throw new ArgumentNullException(nameof(str));

            var strReader = new StringReader(str);
            return LemmatizeHTML(strReader, lem);
        }

        public static IEnumerable<Token> LemmatizeHTML(TextReader reader, Lemmatizer lem)
        {
            if (reader == null) throw new ArgumentNullException(nameof(reader));

            var tokens = TokenizeHTML(reader);
            return Lemmatize(tokens, lem);
        }
    }
}
