using System;
using System.Net;

namespace Jargon.Cmd
{
    internal sealed class CompressedWebClient: WebClient
    {
        protected override WebRequest GetWebRequest(Uri address)
        {
            var @base = base.GetWebRequest(address);
            if(@base is HttpWebRequest http)
            {
                // be a good citizen
                http.AutomaticDecompression = DecompressionMethods.Deflate | DecompressionMethods.GZip;
                // indicate who we are
                http.UserAgent = "Jargon";
                // people _still_ screw this HTTP 1.1 thing up, sigh
                http.Pipelined = false;
            }

            return @base;
        }
    }
}
